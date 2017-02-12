package main

import (
	"encoding/json"
	"net/http"
)

func Groups(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	g := Group{
		"TestGroup",
		"downloads",
		"examplefeed.com",
		nil,
		nil,
		nil,
		nil,
	}

	res := Response{Status{200, ""}, g}
	json.NewEncoder(w).Encode(res)
}

func GroupsCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Body == nil {
		res := Response{Status{400, "No body"}, nil}
		json.NewEncoder(w).Encode(res)
		return
	}
	var g Group
	json.NewDecoder(r.Body).Decode(&g)

	res := Response{Status{200, ""}, nil}
	json.NewEncoder(w).Encode(res)
}

func GroupsCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
}

func GroupsNotifications(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
}

func GroupsShow(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
}

func GroupsDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
}
