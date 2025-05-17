package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"

	_ "github.com/marcboeker/go-duckdb"
)

var db *sql.DB

func initDB() error {
	var err error
	db, err = sql.Open("duckdb", "../db/testdb.ddb")
	if err != nil {
		return fmt.Errorf("error opening db: %w", err)
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS names (name VARCHAR)")
	if err != nil {
		return fmt.Errorf("error creating table: %w", err)
	}
	return nil
}

func addNameToDB(name string) error {
	_, err := db.Exec("INSERT INTO names (name) VALUES (?)", name)
	if err != nil {
		return fmt.Errorf("error inserting name: %w", err)
	}
	return nil
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
		if err := addNameToDB(data.FirstName); err != nil {
			log.Println("Error adding name to DB:", err)
			http.Error(w, "Error adding name to DB", http.StatusInternalServerError)
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
	if err := initDB(); err != nil {
		log.Fatal("Error initializing database:", err)
	}
	defer db.Close()

	http.Handle("/", http.FileServer(http.Dir("../frontend")))
	http.HandleFunc("/api/firstName", enableCORS(fetchData))

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	fmt.Println("Server started at http://localhost:8080")
}
