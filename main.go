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

	groups, err := allGroups()
	if err != nil {
		panic("Failed to init checker with groups.")
	}
	InitChecker(groups)
	defer master.Stop()

	router := NewRouter()
	methods := []string{"GET", "POST", "DELETE", "PUT"}
	log.Fatal(http.ListenAndServe(":8080", handlers.CORS(handlers.AllowedMethods(methods))(router)))
}
