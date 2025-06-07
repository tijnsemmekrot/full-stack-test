package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

var (
	client     *mongo.Client
	collection *mongo.Collection
	err        error
)

func initDB() {
	mongo_pass := os.Getenv("MONGO_DB_PASSWORD")
	uri := "mongodb+srv://tsemmekrot:" + mongo_pass +
		"@full-stack-test.lf9w6dv.mongodb.net/" +
		"?retryWrites=true&w=majority&appName=full-stack-test"
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(opts)
	if err != nil {
		panic(err)
	}

	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	log.Println("Pinged your deployment. You successfully connected to MongoDB!")
	type Person struct {
		Name string
	}
	collection = client.Database("full-stack-test").Collection("names")
	log.Println("Connected to mongoDB")
}

func Handler(w http.ResponseWriter, r *http.Request) {
	type FirstNameRequest struct {
		FirstName string `json:"first_name"`
	}
	type Person struct {
		Name string `bson:"name"`
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}
	var req FirstNameRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Printf("Inserting name: %q into MongoDB", req.FirstName)
	person := Person{Name: req.FirstName}
	result, err := collection.InsertOne(ctx, person)
	if err != nil {
		http.Error(w, "Failed to insert document", http.StatusInternalServerError)
		log.Printf("Insert error: %v\n", err)
		return
	}

	var idStr string
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		idStr = oid.Hex()
	} else {
		idStr = fmt.Sprintf("%v", result.InsertedID)
	}
	log.Printf("Inserted document with ID: %v\n", idStr)

	type Response struct {
		Message string `json:"message"`
		ID      string `json:"id"`
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{Message: req.FirstName + " added to MongoDB!", ID: idStr})
}

func getData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		http.Error(w, "Failed to retrieve documents", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	type person struct {
		ID   string `json:"_id"`
		Name string `bson:"name" json:"name"`
	}

	var persons []person

	// Iterate through each document
	for cursor.Next(ctx) {
		var doc struct {
			ID   interface{} `bson:"_id"`
			Name string      `bson:"name"`
		}

		if err := cursor.Decode(&doc); err != nil {
			http.Error(w, "Failed to decode document", http.StatusInternalServerError)
			return
		}

		// Convert ID to string
		var idStr string
		if oid, ok := doc.ID.(primitive.ObjectID); ok {
			idStr = oid.Hex()
		} else {
			idStr = fmt.Sprintf("%v", doc.ID)
		}

		persons = append(persons, person{
			ID:   idStr,
			Name: doc.Name,
		})
	}

	if err := cursor.Err(); err != nil {
		http.Error(w, "Cursor error", http.StatusInternalServerError)
		return
	}

	log.Printf("Retrieved %d documents\n", len(persons))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(persons)
}

// test
func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func main() {
	goVersion := os.Getenv("GO_VERSION")
	log.Println("GO_VERSION:", goVersion)

	initDB()

	http.HandleFunc("/api/firstName", enableCORS(Handler))
	http.HandleFunc("/api/getData", enableCORS(getData))

	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
