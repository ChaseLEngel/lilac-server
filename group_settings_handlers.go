package main

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
)

func Settings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var res Response
	group, err := findGroup(mux.Vars(r)["groupId"])
	if err != nil {
		res = NewResponse(400, err, nil)
		json.NewEncoder(w).Encode(res)
		return
	}
	settings, err := group.GroupSettings()
	if err != nil {
		res = NewResponse(500, err, nil)
		json.NewEncoder(w).Encode(res)
		return
	}
	res = NewResponse(200, nil, settings)
	json.NewEncoder(w).Encode(res)
}

func SettingsCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var res Response
	group, err := findGroup(mux.Vars(r)["groupId"])
	if err != nil {
		res = NewResponse(400, err, nil)
		json.NewEncoder(w).Encode(res)
		return
	}

	if r.Body == nil {
		res = NewResponse(400, errors.New("No body"), nil)
		json.NewEncoder(w).Encode(res)
		return
	}

	var settings GroupSettings
	json.NewDecoder(r.Body).Decode(&settings)

	err = group.insertGroupSettings(settings)
	if err != nil {
		res = NewResponse(500, err, nil)
		json.NewEncoder(w).Encode(res)
		return
	}
	res = NewResponse(200, nil, settings)
	json.NewEncoder(w).Encode(res)
}
