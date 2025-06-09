package db

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	Client     *mongo.Client
	Collection *mongo.Collection
)

func InitDB() {
	mongoPass := os.Getenv("MONGO_DB_PASSWORD")
	uri := "mongodb+srv://tsemmekrot:" + mongoPass +
		"@full-stack-test.lf9w6dv.mongodb.net/" +
		"?retryWrites=true&w=majority&appName=full-stack-test"
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	var err error
	Client, err = mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := Client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}
	log.Println("Pinged your deployment. You successfully connected to MongoDB!")

	Collection = Client.Database("full-stack-test").Collection("names")
	log.Println("Connected to MongoDB")
}
