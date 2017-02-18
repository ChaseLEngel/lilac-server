package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func GroupsConstraints(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var res Response
	group, err := findGroup(mux.Vars(r)["groupId"])
	if err != nil {
		res = Response{Status{400, err.Error()}, nil}
		json.NewEncoder(w).Encode(res)
		return
	}
	constraints, err := group.allConstraints()
	if err != nil {
		res = Response{Status{400, err.Error()}, nil}
		json.NewEncoder(w).Encode(res)
		return
	}
	res = Response{Status{200, ""}, constraints}
	json.NewEncoder(w).Encode(res)

}

func GroupsConstraintsCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var res Response
	if r.Body == nil {
		res = Response{Status{400, "No body"}, nil}
		json.NewEncoder(w).Encode(res)
		return
	}
	group, err := findGroup(mux.Vars(r)["groupId"])
	if err != nil {
		res = Response{Status{400, err.Error()}, nil}
		json.NewEncoder(w).Encode(res)
		return
	}
	var constraint Constraint
	json.NewDecoder(r.Body).Decode(&constraint)
	err = group.insertConstraint(&constraint)
	if err != nil {
		res = Response{Status{400, err.Error()}, nil}
		json.NewEncoder(w).Encode(res)
		return
	}
	res = Response{Status{200, ""}, constraint}
	json.NewEncoder(w).Encode(res)
}

func GroupsConstraintsShow(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var res Response
	group, err := findGroup(mux.Vars(r)["groupId"])
	if err != nil {
		res = Response{Status{400, err.Error()}, nil}
		json.NewEncoder(w).Encode(res)
		return
	}
	constraint, err := group.findConstraint(mux.Vars(r)["constraintId"])
	if err != nil {
		res = Response{Status{400, err.Error()}, nil}
		json.NewEncoder(w).Encode(res)
		return
	}
	res = Response{Status{200, ""}, constraint}
	json.NewEncoder(w).Encode(res)
}

func GroupsConstraintsDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var res Response
	group, err := findGroup(mux.Vars(r)["groupId"])
	if err != nil {
		res = Response{Status{400, err.Error()}, nil}
		json.NewEncoder(w).Encode(res)
		return
	}
	err = group.deleteConstraint(mux.Vars(r)["constraintId"])
	if err != nil {
		res = Response{Status{400, err.Error()}, nil}
		json.NewEncoder(w).Encode(res)
		return
	}
	res = Response{Status{200, ""}, nil}
	json.NewEncoder(w).Encode(res)
}

func GroupsConstraintsUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var res Response

	if r.Body == nil {
		res = Response{Status{400, "No body"}, nil}
		json.NewEncoder(w).Encode(res)
		return
	}

	group, err := findGroup(mux.Vars(r)["groupId"])
	if err != nil {
		res = Response{Status{400, err.Error()}, nil}
		json.NewEncoder(w).Encode(res)
		return
	}

	var constraint Constraint
	json.NewDecoder(r.Body).Decode(&constraint)

	err = group.updateConstraint(mux.Vars(r)["constraintId"], &constraint)
	if err != nil {
		res = Response{Status{400, err.Error()}, nil}
		json.NewEncoder(w).Encode(res)
		return
	}

	res = Response{Status{200, ""}, constraint}
	json.NewEncoder(w).Encode(res)
}
