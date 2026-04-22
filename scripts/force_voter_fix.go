package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	mongoURI := "mongodb+srv://Admin:Cv3hzz3fta@database.qnmktuh.mongodb.net/?appName=Database"
	dbName := "college_data"
	email := "harshitpandeyblp2021@gmail.com"
	address := "0x16Fc9d5e4FF4a55B1d43a01016C190C607848869"

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	db := client.Database(dbName)
	
	// FIX: Update status to "Voted" using case-insensitive filter
	// This fixes the user's current session
	res, err := db.Collection("voters").UpdateOne(context.Background(),
		bson.M{
			"email": bson.M{"$regex": "^" + strings.TrimSpace(email) + "$", "$options": "i"},
			"registrations.election_address": bson.M{"$regex": "^" + strings.TrimSpace(address) + "$", "$options": "i"},
		},
		bson.M{"$set": bson.M{"registrations.$.status": "Voted"}},
	)

	if err != nil {
		fmt.Println("Error updating voter:", err)
	} else {
		fmt.Printf("Successfully updated %d voter document\n", res.ModifiedCount)
	}
}
