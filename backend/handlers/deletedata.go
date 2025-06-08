package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"full-stack-test/db"
	"full-stack-test/models"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo-driver/v2/primitive"
)

func DeleteData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.IdDeleteRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	if req.Id == "" {
		http.Error(w, "ID cannot be empty", http.StatusBadRequest)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid ID format: %v", err), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := db.Collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		http.Error(w, fmt.Sprintf("Error deleting document: %v", err), http.StatusInternalServerError)
		return
	}

	if result.DeletedCount == 0 {
		http.Error(w, "No document found with the specified ID", http.StatusNotFound)
		return
	}

	response := models.DeleteResponse{
		Message: req.Id + " removed from MongoDB!",
		Id:      req.Id,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
	log.Printf("Deleting document with ID: %v", req.Id)
}
