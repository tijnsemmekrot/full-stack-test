package main

import (
	"context"
	//"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

func initDB() {
	MONGO_PASSWORD := os.Getenv("MONGO_DB_PASSWORD")
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI("mongodb+srv://tsemmekrot:" + MONGO_PASSWORD + "@full-stack-test.lf9w6dv.mongodb.net/?retryWrites=true&w=majority&appName=full-stack-test").SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(opts)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	// Send a ping to confirm a successful connection
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
}

//func addNameToDB(name string) error {
//	_, err := conn.Exec(context.Background(),
//		"INSERT INTO names (name) VALUES ($1)",
//		name,
//	)
//	if err != nil {
//		return fmt.Errorf("error inserting name: %w", err)
//	}
//	return nil
//}

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

//func fetchData(w http.ResponseWriter, r *http.Request) {
//	if r.Method == "POST" {
//		var data struct {
//			FirstName string `json:"first_name"`
//		}
//		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
//			http.Error(w, "Invalid JSON", http.StatusBadRequest)
//			return
//		}
//		if err := addNameToDB(data.FirstName); err != nil {
//			log.Println("Error adding name to DB:", err)
//			http.Error(w, "Error adding name to DB", http.StatusInternalServerError)
//			return
//		}
//		log.Println("Request received with first name:", data.FirstName)
//		w.Header().Set("Content-Type", "application/json")
//		json.NewEncoder(w).Encode(map[string]string{
//			"status":   "success",
//			"received": data.FirstName,
//		})
//	}
//}

func main() {
	log.Println("Go version:", runtime.Version())
	goVersion := os.Getenv("GO_VERSION")
	log.Println("GO_VERSION:", goVersion)
	initDB()

	// http.Handle("/", http.FileServer(http.Dir("../frontend")))
	// http.HandleFunc("/api/firstName", enableCORS(fetchData))
	//
	// err := http.ListenAndServe(":8080", nil)
	//
	//	if err != nil {
	//		log.Fatal("ListenAndServe: ", err)
	//	}
	//
	// fmt.Println("Server started at http://localhost:8080")
}
