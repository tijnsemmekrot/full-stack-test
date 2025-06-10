package handlers

import (
	"context"
	"encoding/json"
	"full-stack-test/models"
	"log"
	"net/http"
	"time"
)

func ExpenseHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var input models.InsertExpenseRequest
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	result, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("results: ", result)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.Response{Message: "Expense inserted successfully", ID: "123"})
}
