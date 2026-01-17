package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	collection := client.Database(dbName).Collection("candidates")

	// Delete all documents
	res, err := collection.DeleteMany(ctx, map[string]interface{}{})
	if err != nil {
		log.Fatalf("Failed to delete candidates: %v", err)
	}

	fmt.Printf("✅ Successfully deleted %d candidate(s) from the database.\n", res.DeletedCount)
}
