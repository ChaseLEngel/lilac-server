package main

import (
	"github.com/gorilla/handlers"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
	"net/http"
)

var db *gorm.DB

func main() {
	router := NewRouter()
	var err error
	db, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("Failed to connect to database.")
	}
	log.Fatal(http.ListenAndServe(":8080", handlers.CORS()(router)))
}
