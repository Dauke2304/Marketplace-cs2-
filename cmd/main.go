package main

import (
	"Marketplace-cs2-/database"
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func main() {
	database.ConnectDB()

	usersCollection := database.GetCollection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var result bson.M
	err := usersCollection.FindOne(ctx, bson.M{"steam_id": "test_steam_id"}).Decode(&result)
	if err != nil {
		log.Fatalf("Error fetching user: %v", err)
	}

	fmt.Printf("Fetched user: %v\n", result)
}
