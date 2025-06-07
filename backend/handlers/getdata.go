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
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetDataHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := db.Collection.Find(ctx, bson.D{})
	if err != nil {
		http.Error(w, "Failed to retrieve documents", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var persons []models.PersonResponse

	for cursor.Next(ctx) {
		var doc struct {
			ID   interface{} `bson:"_id"`
			Name string      `bson:"name"`
		}

		if err := cursor.Decode(&doc); err != nil {
			http.Error(w, "Failed to decode document", http.StatusInternalServerError)
			return
		}

		var idStr string
		if oid, ok := doc.ID.(primitive.ObjectID); ok {
			idStr = oid.Hex()
		} else {
			idStr = fmt.Sprintf("%v", doc.ID)
		}

		persons = append(persons, models.PersonResponse{
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
