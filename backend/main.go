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

	http.HandleFunc("/api/firstName", middleware.EnableCORS(handlers.Handler))
	http.HandleFunc("/api/getData", middleware.EnableCORS(handlers.GetDataHandler))

	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
