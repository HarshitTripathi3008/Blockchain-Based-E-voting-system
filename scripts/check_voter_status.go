package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	mongoURI := "mongodb+srv://Admin:Cv3hzz3fta@database.qnmktuh.mongodb.net/?appName=Database"
	dbName := "college_data"

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	db := client.Database(dbName)
	
	// Check ALL Jobs
	cursor, err := db.Collection("vote_jobs").Find(context.Background(), bson.M{})
	fmt.Println("\n--- ALL VOTE JOBS (Limited to 20) ---")
	if err != nil {
		fmt.Println("Error finding jobs:", err)
	} else {
		var jobs []bson.M
		cursor.All(context.Background(), &jobs)
		if len(jobs) > 20 {
			jobs = jobs[:20]
		}
		b, _ := json.MarshalIndent(jobs, "", "  ")
		fmt.Println(string(b))
	}
}
