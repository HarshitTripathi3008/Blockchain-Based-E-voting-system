package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
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
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("Note: No .env file found in parent dir, checking current...")
		if err := godotenv.Load(".env"); err != nil {
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

	// 1. Delete Cloudinary Images
	fmt.Println("🔍 Fetching candidates to clear images...")
	cursor, err := candColl.Find(ctx, bson.M{})
	if err == nil {
		defer cursor.Close(ctx)
		var candidates []bson.M
		if err := cursor.All(ctx, &candidates); err == nil {
			for _, c := range candidates {
				if imgHash, ok := c["imageHash"].(string); ok && imgHash != "" {
					// Extract public ID logic
					// Expected URL: .../upload/v123/voting_system/filename.png
					// We need: voting_system/filename (no ext)

					// Simple hack: find "voting_system" and take rest
					if idx := strings.Index(imgHash, "voting_system"); idx != -1 {
						part := imgHash[idx:] // voting_system/filename.png
						ext := filepath.Ext(part)
						publicID := strings.TrimSuffix(part, ext)

						fmt.Printf("🗑️ Deleting image: %s ... ", publicID)
						if err := util.DeleteFromCloudinary(publicID); err != nil {
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
		fmt.Printf("✅ Deleted %d candidates.\n", res.DeletedCount)
	}

	// 3. Delete OTPs
	resOtp, err := otpColl.DeleteMany(ctx, bson.M{})
	if err != nil {
		log.Printf("Failed to delete OTPs: %v", err)
	} else {
		fmt.Printf("✅ Deleted %d OTPs.\n", resOtp.DeletedCount)
	}

	// 4. Delete Audit Logs
	resAudit, err := auditColl.DeleteMany(ctx, bson.M{})
	if err != nil {
		log.Printf("Failed to delete Audit Logs: %v", err)
	} else {
		fmt.Printf("✅ Deleted %d Audit Log entries.\n", resAudit.DeletedCount)
	}
}
