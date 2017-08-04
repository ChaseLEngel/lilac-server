package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func Groups(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	groups, err := allGroups()
	var res Response
	if err != nil {
		res = NewResponse(400, err, nil)
	} else {
		res = NewResponse(200, nil, groups)
	}
	json.NewEncoder(w).Encode(res)
}

func GroupsCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var res Response
	if r.Body == nil {
		res = NewResponse(400, fmt.Errorf("No body"), nil)
		json.NewEncoder(w).Encode(res)
		return
	}

	var group Group
	json.NewDecoder(r.Body).Decode(&group)

	err := group.insert()
	if err != nil {
		res = NewResponse(400, err, nil)
		json.NewEncoder(w).Encode(res)
		return
	}

	res = NewResponse(200, nil, group)
	json.NewEncoder(w).Encode(res)
}

// Run check on a single group.
func GroupsCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var res Response
	group, err := findGroup(mux.Vars(r)["groupId"])
	if err != nil {
		res = NewResponse(500, err, nil)
		json.NewEncoder(w).Encode(res)
		return
	}
	check(group.ID)
	res = NewResponse(200, nil, nil)
	json.NewEncoder(w).Encode(res)
}

func GroupsShow(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var res Response
	group, err := findGroup(mux.Vars(r)["groupId"])
	if err != nil {
		res = NewResponse(400, err, nil)
	} else {
		res = NewResponse(200, nil, group)
	}
	json.NewEncoder(w).Encode(res)
}

func GroupsDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var res Response
	group, err := findGroup(mux.Vars(r)["groupId"])
	if err != nil {
		res = NewResponse(400, err, nil)
		json.NewEncoder(w).Encode(res)
		return
	}
	err = deleteGroup(mux.Vars(r)["groupId"])
	if err != nil {
		res = NewResponse(400, err, nil)
		json.NewEncoder(w).Encode(res)
		return
	}
	res = NewResponse(200, nil, group)
	json.NewEncoder(w).Encode(res)
}

func GroupsUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var res Response

	if r.Body == nil {
		res = NewResponse(400, fmt.Errorf("No body"), nil)
		json.NewEncoder(w).Encode(res)
		return
	}

	group, err := findGroup(mux.Vars(r)["groupId"])
	if err != nil {
		res = NewResponse(400, err, nil)
		json.NewEncoder(w).Encode(res)
		return
	}

	var updateGroup Group
	err = json.NewDecoder(r.Body).Decode(&updateGroup)
	if err != nil {
		res = NewResponse(400, err, nil)
		json.NewEncoder(w).Encode(res)
		return
	}

	err = group.update(updateGroup)
	if err != nil {
		res = NewResponse(400, err, nil)
		json.NewEncoder(w).Encode(res)
		return
	}

	res = NewResponse(200, nil, group)
	json.NewEncoder(w).Encode(res)
}
