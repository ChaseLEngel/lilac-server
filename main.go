package main

import (
	"github.com/gorilla/handlers"
	"log"
	"net/http"
)

func main() {
	router := NewRouter()
	log.Fatal(http.ListenAndServe(":8080", handlers.CORS()(router)))
}
