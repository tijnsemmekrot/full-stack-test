package main

import (
	"context"
	"log"
	"net/url"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

func initDB() {
	// Encode special characters in password
	rawPassword := os.Getenv("MONGO_DB_PASSWORD")
	MONGO_PASSWORD := url.QueryEscape(rawPassword)

	uri := "mongodb+srv://tsemmekrot:" + MONGO_PASSWORD +
		"@full-stack-test.lf9w6dv.mongodb.net/" +
		"?retryWrites=true&w=majority&appName=full-stack-test"

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().
		ApplyURI(uri).
		SetServerAPIOptions(serverAPI)

	// Add timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		log.Fatal("Connection failed:", err)
	}

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Fatal("Disconnect failed:", err)
		}
	}()

	// Verify connection
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal("Ping failed:", err)
	}
	log.Println("Successfully connected to MongoDB!")
}

func main() {
	initDB()
}
