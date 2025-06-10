package main

import (
	"full-stack-test/db"
	"full-stack-test/handlers"
	"full-stack-test/middleware"
	"log"
	"net/http"
	"os"
)

func main() {
	goVersion := os.Getenv("GO_VERSION")
	log.Println("GO_VERSION:", goVersion)

	db.InitDB()

	http.HandleFunc("/api/firstName", middleware.EnableCORS(handlers.InsertFirstName))
	http.HandleFunc("/api/getData", middleware.EnableCORS(handlers.GetDataHandler))
	http.HandleFunc("/api/deleteData", middleware.EnableCORS(handlers.DeleteData))
	http.HandleFunc("/api/expense", middleware.EnableCORS(handlers.ExpenseHandler))
	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
