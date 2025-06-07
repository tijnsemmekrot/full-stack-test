// package handlers
//
//import (
//	"context"
//	"encoding/json"
//	"fmt"
//	"full-stack-test/db"
//	"full-stack-test/models"
//	"log"
//	"net/http"
//	"time"
//
//	"go.mongodb.org/mongo-driver/bson/primitive"
//)
//
//func DeleteFirstName(w http.ResponseWriter, r *http.Request) {
//	if r.Method != http.MethodPost {
//		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
//		return
//	}
//
//	var req models.PersonResponse
//	err := json.NewDecoder(r.Body).Decode(&req)
//	if err != nil {
//		http.Error(w, "Invalid request body", http.StatusBadRequest)
//		return
//	}
//
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//
//	log.Printf("Inserting name: %q into MongoDB", req.FirstName)
//	person := models.Person{Name: req.FirstName}
//	result, err := db.Collection.InsertOne(ctx, person)
//	if err != nil {
//		http.Error(w, "Failed to insert document", http.StatusInternalServerError)
//		log.Printf("Insert error: %v\n", err)
//		return
//	}
//
//	var idStr string
//	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
//		idStr = oid.Hex()
//	} else {
//		idStr = fmt.Sprintf("%v", result.InsertedID)
//	}
//	log.Printf("Inserted document with ID: %v\n", idStr)
//
//	response := models.Response{
//		Message: req.FirstName + " added to MongoDB!",
//		ID:      idStr,
//	}
//	w.Header().Set("Content-Type", "application/json")
//	w.WriteHeader(http.StatusOK)
//	json.NewEncoder(w).Encode(response)
//}
