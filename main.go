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
	log.Fatal(http.ListenAndServe(":8080", handlers.CORS()(router)))
}
