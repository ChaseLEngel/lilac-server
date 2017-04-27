package main

import (
	"github.com/chaselengel/lilac/logger"
	"github.com/gorilla/handlers"
	"net/http"
)

var log *logger.Logger

func main() {
	log = logger.New("./lilac.log")

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
	http.ListenAndServe(":8080", handlers.CORS(handlers.AllowedMethods(methods))(router))
}
