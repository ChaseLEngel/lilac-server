package main

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
)

func Requests(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	group, err := findGroup(mux.Vars(r)["groupId"])
	var res Response
	if err != nil {
		res = NewResponse(400, err, nil)
		json.NewEncoder(w).Encode(res)
		return
	}
	requests, err := group.allRequests()
	if err != nil {
		res = NewResponse(500, err, nil)
		json.NewEncoder(w).Encode(res)
		return
	}
	res = NewResponse(200, nil, requests)
	json.NewEncoder(w).Encode(res)
}

func RequestsCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Body == nil {
		res := NewResponse(400, errors.New("No body"), nil)
		json.NewEncoder(w).Encode(res)
		return
	}
	group, err := findGroup(mux.Vars(r)["groupId"])
	if err != nil {
		res := NewResponse(400, err, nil)
		json.NewEncoder(w).Encode(res)
		return
	}
	var request Request
	json.NewDecoder(r.Body).Decode(&request)

	if err := group.insertRequest(&request); err != nil {
		res := NewResponse(500, err, nil)
		json.NewEncoder(w).Encode(res)
		return
	}

	res := NewResponse(200, nil, request)
	json.NewEncoder(w).Encode(res)
}

func RequestsDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	group, err := findGroup(mux.Vars(r)["groupId"])
	var res Response
	if err != nil {
		res = NewResponse(400, err, nil)
		json.NewEncoder(w).Encode(res)
		return
	}
	request, err := group.findRequest(mux.Vars(r)["requestId"])
	if err != nil {
		res = NewResponse(400, err, nil)
		json.NewEncoder(w).Encode(res)
		return
	}
	err = group.deleteRequest(mux.Vars(r)["requestId"])
	if err != nil {
		res = NewResponse(400, err, nil)
		json.NewEncoder(w).Encode(res)
		return
	}
	res = NewResponse(200, nil, request)
	json.NewEncoder(w).Encode(res)
}

func RequestsUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	group, err := findGroup(mux.Vars(r)["groupId"])
	var res Response
	if err != nil {
		res = NewResponse(400, err, nil)
		json.NewEncoder(w).Encode(res)
		return
	}

	var request Request
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		res = NewResponse(400, err, nil)
		json.NewEncoder(w).Encode(res)
		return
	}

	updatedRequest, err := group.updateRequest(mux.Vars(r)["requestId"], request)
	if err != nil {
		res = NewResponse(400, err, nil)
		json.NewEncoder(w).Encode(res)
		return
	}

	res = NewResponse(200, nil, updatedRequest)
	json.NewEncoder(w).Encode(res)
}

func RequestsHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var res Response
	group, err := findGroup(mux.Vars(r)["groupId"])
	if err != nil {
		res = NewResponse(400, err, nil)
		json.NewEncoder(w).Encode(res)
		return
	}
	request, err := group.findRequest(mux.Vars(r)["requestId"])
	if err != nil {
		res = NewResponse(500, err, nil)
		json.NewEncoder(w).Encode(res)
		return
	}
	matchHistory, err := request.history()
	if err != nil {
		res = NewResponse(500, err, nil)
		json.NewEncoder(w).Encode(res)
		return
	}
	res = NewResponse(200, nil, matchHistory)
	json.NewEncoder(w).Encode(res)
}

func RequestsHistoryDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var res Response
	group, err := findGroup(mux.Vars(r)["groupId"])
	if err != nil {
		res = NewResponse(400, err, nil)
		json.NewEncoder(w).Encode(res)
		return
	}
	request, err := group.findRequest(mux.Vars(r)["requestId"])
	if err != nil {
		res = NewResponse(500, err, nil)
		json.NewEncoder(w).Encode(res)
		return
	}

	history, err := request.deleteMatchHistory(mux.Vars(r)["historyId"])
	if err != nil {
		res = NewResponse(500, err, nil)
		json.NewEncoder(w).Encode(res)
		return
	}

	res = NewResponse(200, nil, history)
	json.NewEncoder(w).Encode(res)
}
