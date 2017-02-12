package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func Groups(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	groups := Db.Find(&([]Group{})).Value
	res := Response{Status{200, ""}, groups}
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

	Db.NewRecord(g)
	Db.Create(&g)

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
	group, err := findGroup(mux.Vars(r)["groupId"])
	res := Response{}
	if err != nil {
		res = Response{Status{400, err.Error()}, nil}
	} else {
		res = Response{Status{200, ""}, group}
	}
	json.NewEncoder(w).Encode(res)
}

func GroupsDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := deleteGroup(mux.Vars(r)["groupId"])
	res := Response{}
	if err != nil {
		res = Response{Status{400, err.Error()}, nil}
	} else {
		res = Response{Status{200, ""}, nil}
	}
	json.NewEncoder(w).Encode(res)
}
