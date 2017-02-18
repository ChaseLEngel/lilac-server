package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func GroupsConstraints(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
}
