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
		json.NewEncoder(w).Encode(res)
		return
	}
	requests, err := group.allRequests()
	if err != nil {
		res = Response{Status{500, err.Error()}, nil}
	} else {
		res = Response{Status{200, ""}, requests}
	}
	json.NewEncoder(w).Encode(res)
}

func RequestsCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Body == nil {
		res := Response{Status{400, "No body"}, nil}
		json.NewEncoder(w).Encode(res)
		return
	}
	var request Request
	group, err := findGroup(mux.Vars(r)["groupId"])
	if err != nil {
		res := Response{Status{400, err.Error()}, nil}
		json.NewEncoder(w).Encode(res)
		return
	}
	json.NewDecoder(r.Body).Decode(&request)

	if err := group.insertRequest(request); err != nil {
		res := Response{Status{500, err.Error()}, nil}
		json.NewEncoder(w).Encode(res)
		return
	}

	res := Response{Status{200, ""}, nil}
	json.NewEncoder(w).Encode(res)
}

func RequestsDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
}

func RequestsHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
}
