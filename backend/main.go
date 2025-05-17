package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var conn *pgx.Conn

func initDB() error {
	readSecret := func(filename string) (string, error) {
		// Try Render's secret path first
		content, err := os.ReadFile("/etc/secrets/" + filename)
		if err != nil {
			// Fallback to local development
			content, err = os.ReadFile(".env." + filename)
			if err != nil {
				return "", fmt.Errorf("missing %s: %w", filename, err)
			}
		}
		return strings.TrimSpace(string(content)), nil
	}

	host, _ := readSecret("DB_HOST")
	port, _ := readSecret("DB_PORT")
	user, _ := readSecret("DB_USER")
	password, _ := readSecret("DB_PASSWORD")
	dbname, _ := readSecret("DB_NAME")
	// Read required secrets

	if host == "" || user == "" || dbname == "" {
		return fmt.Errorf("missing required database credentials")
	}

	connConfig, err := pgconn.ParseConfig(fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		host, port, user, password, dbname,
	))
	if err != nil {
		return fmt.Errorf("config parse failed: %w", err)
	}

	connConfig.RuntimeParams["auth_type"] = "scram-sha-256"

	connConfig.TLSConfig = &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	log.Printf("Connecting to: postgres://%s:***@%s:%s/%s", user, host, port, dbname)

	conn, err = pgx.Connect(context.Background(), connConfig)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %w", err)
	}
	if err := conn.Ping(context.Background()); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}
	_, err = conn.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS names (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT NOW()
		)
	`)
	if err != nil {
		return fmt.Errorf("error creating table: %w", err)
	}
	return nil
}

func addNameToDB(name string) error {
	_, err := conn.Exec(context.Background(),
		"INSERT INTO names (name) VALUES ($1)",
		name,
	)
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
	goVersion := os.Getenv("GO_VERSION")
	log.Println("GO_VERSION:", goVersion)
	if err := initDB(); err != nil {
		log.Fatal("Error initializing database:", err)
	}
	defer conn.Close(context.Background())

	http.Handle("/", http.FileServer(http.Dir("../frontend")))
	http.HandleFunc("/api/firstName", enableCORS(fetchData))

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	fmt.Println("Server started at http://localhost:8080")
}
