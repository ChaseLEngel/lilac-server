package main

import (
	"github.com/gorilla/handlers"
	"log"
	"net/http"
)

func main() {
	err := initDatabase("production.db")
	if err != nil {
		panic("Failed to init database")
	}
	router := NewRouter()
	methods := []string{"GET", "POST", "DELETE"}
	log.Fatal(http.ListenAndServe(":8080", handlers.CORS(handlers.AllowedMethods(methods))(router)))
}
