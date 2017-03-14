package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func Settings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var res Response
	group, err := findGroup(mux.Vars(r)["groupId"])
	if err != nil {
		res = Response{Status{400, err.Error()}, nil}
		json.NewEncoder(w).Encode(res)
		return
	}
	settings, err := group.GroupSettings()
	if err != nil {
		res = Response{Status{500, err.Error()}, nil}
		json.NewEncoder(w).Encode(res)
		return
	}
	res = Response{Status{200, ""}, settings}
	json.NewEncoder(w).Encode(res)
}

func SettingsCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var res Response
	group, err := findGroup(mux.Vars(r)["groupId"])
	if err != nil {
		res = Response{Status{400, err.Error()}, nil}
		json.NewEncoder(w).Encode(res)
		return
	}

	if r.Body == nil {
		res = Response{Status{400, "No Body"}, nil}
		json.NewEncoder(w).Encode(res)
		return
	}

	var settings GroupSettings
	json.NewDecoder(r.Body).Decode(&settings)

	err = group.insertGroupSettings(settings)
	if err != nil {
		res = Response{Status{500, err.Error()}, nil}
		json.NewEncoder(w).Encode(res)
		return
	}
	res = Response{Status{200, ""}, settings}
	json.NewEncoder(w).Encode(res)
}
