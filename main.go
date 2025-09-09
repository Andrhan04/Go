package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	if err := InitDB(); err != nil {
		log.Fatal("Error initializing database:", err)
	}
	defer CloseDB()
	r := mux.NewRouter()

	r.HandleFunc("/cats", CreateCatHandler).Methods("POST")
	r.HandleFunc("/cats", GetCatsHandler).Methods("GET")
	r.HandleFunc("/cats/{id}", GetCatHandler).Methods("GET")
	r.HandleFunc("/cats/{id}", UpdateCatHandler).Methods("PUT")
	r.HandleFunc("/cats/{id}", DeleteCatHandler).Methods("DELETE")

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
