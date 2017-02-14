package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func Requests(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	group, err := findGroup(mux.Vars(r)["groupId"])
	var res Response
	if err != nil {
		res = Response{Status{400, err.Error()}, nil}
		return
	}
	requests, err := group.allRequests()
	if err != nil {
		res = Response{Status{400, err.Error()}, nil}
	} else {
		res = Response{Status{200, ""}, requests}
	}
	json.NewEncoder(w).Encode(res)
}

func RequestsCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
}

func RequestsDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
}

func RequestsHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
}
