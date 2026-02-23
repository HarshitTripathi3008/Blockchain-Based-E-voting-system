package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"MAJOR-PROJECT/util"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Load .env
	if err := godotenv.Overload("../.env"); err != nil {
		log.Println("Note: No .env file found in parent dir, checking current...")
		if err := godotenv.Overload(".env"); err != nil {
			log.Println("Note: No .env file found, relying on system env")
		}
	}

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "voting_system"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second) // Increased timeout
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	db := client.Database(dbName)
	candColl := db.Collection("candidates")
	otpColl := db.Collection("otps")
	auditColl := db.Collection("audit_logs")

	// 1. Delete S3 Images
	fmt.Println("[SEARCH] Fetching candidates to clear images from S3...")
	cursor, err := candColl.Find(ctx, bson.M{})
	if err == nil {
		defer cursor.Close(ctx)
		var candidates []bson.M
		if err := cursor.All(ctx, &candidates); err == nil {
			for _, c := range candidates {
				if imgUrl, ok := c["imageHash"].(string); ok && imgUrl != "" {
					// Extract S3 key logic
					// Expected URL: https://bucket.s3.region.amazonaws.com/uploads/timestamp_filename
					// Key: uploads/timestamp_filename

					if idx := strings.Index(imgUrl, "/uploads/"); idx != -1 {
						objectKey := imgUrl[idx+1:] // uploads/timestamp_filename

						fmt.Printf("[DELETE] Deleting S3 object: %s ... ", objectKey)
						if err := util.DeleteFromS3(objectKey); err != nil {
							fmt.Printf("Failed: %v\n", err)
						} else {
							fmt.Printf("Done\n")
						}
					}
				}
			}
		}
	}

	// 2. Delete Candidates
	res, err := candColl.DeleteMany(ctx, bson.M{})
	if err != nil {
		log.Printf("Failed to delete candidates: %v", err)
	} else {
		fmt.Printf("[OK] Deleted %d candidates.\n", res.DeletedCount)
	}

	// 3. Delete OTPs
	resOtp, err := otpColl.DeleteMany(ctx, bson.M{})
	if err != nil {
		log.Printf("Failed to delete OTPs: %v", err)
	} else {
		fmt.Printf("[OK] Deleted %d OTPs.\n", resOtp.DeletedCount)
	}

	// 4. Delete Audit Logs
	resAudit, err := auditColl.DeleteMany(ctx, bson.M{})
	if err != nil {
		log.Printf("Failed to delete Audit Logs: %v", err)
	} else {
		fmt.Printf("[OK] Deleted %d Audit Log entries.\n", resAudit.DeletedCount)
	}

	// 5. Delete Election Metadata
	metadataColl := db.Collection("election_metadata")
	resMeta, err := metadataColl.DeleteMany(ctx, bson.M{})
	if err != nil {
		log.Printf("Failed to delete Election Metadata: %v", err)
	} else {
		fmt.Printf("[OK] Deleted %d Election Metadata entries.\n", resMeta.DeletedCount)
	}
}
