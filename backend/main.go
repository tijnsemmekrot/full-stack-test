package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"

	"github.com/marcboeker/go-duckdb"
)

func connectToDatabase() {
  db, err := sql.Open("duckdb", "../db/testdb.ddb")
  if err != nil {
    log.Println("error opening db:", err))
  }
}

func addNameToDB(name String) {
  db, err := sql.Open("duckdb", "../db/testdb.ddb")
  if err != nil {
    log.Println("error opening db:", err))
  }
  db.Exec("CREATE TABLE IF NOT EXISTS names (name VARCHAR)")
  _, err = db.Exec("INSERT INTO names (name) VALUES (?)", name)
  if err != nil {
    log.Println("error inserting name:", err)
  }
}


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

func fetchData(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var data struct {
			FirstName string `json:"first_name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		log.Println("Request received with first name:", data.FirstName)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":   "success",
			"received": data.FirstName,
		})
	}
}

func main() {
	log.Println("Go version:", runtime.Version())
  connectToDatabase()
	http.Handle("/", http.FileServer(http.Dir("../frontend")))
  firstName := http.HandleFunc("/api/firstName", enableCORS(fetchData))
  addNameToDB(firstName)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	fmt.Println("Server started at http://localhost:8080")
}
