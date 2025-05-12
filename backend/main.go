package main

import (
	"fmt"
	"log"
	"net/http"
)

func fetchData() {
	http.HandleFunc("/api/firstName", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			firstName := r.FormValue("first_name")
			fmt.Printf("Request received with first name: %s", firstName)
			w.Write([]byte("First name received"))
		}
	})
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("../frontend")))
	fetchData()
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	fmt.Println("Server started at http://localhost:8080")
}
