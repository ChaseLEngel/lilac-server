package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func Groups(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	groups, err := allGroups()
	var res Response
	if err != nil {
		res = Response{Status{400, err.Error()}, nil}
	} else {
		res = Response{Status{200, ""}, groups}
	}
	json.NewEncoder(w).Encode(res)
}

func GroupsCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Body == nil {
		res := Response{Status{400, "No body"}, nil}
		json.NewEncoder(w).Encode(res)
		return
	}
	var group Group
	json.NewDecoder(r.Body).Decode(&group)

	err := insertGroup(&group)
	if err != nil {
		res := Response{Status{500, err.Error()}, nil}
		json.NewEncoder(w).Encode(res)
		return
	}

	res := Response{Status{200, ""}, group}
	json.NewEncoder(w).Encode(res)
}

// Run check on a single group.
func GroupsCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var res Response
	group, err := findGroup(mux.Vars(r)["groupId"])
	if err != nil {
		res = Response{Status{500, err.Error()}, nil}
		json.NewEncoder(w).Encode(res)
		return
	}
	var groups []Group
	groups = append(groups, group)
	go check(groups)
	res = Response{Status{200, ""}, nil}
	json.NewEncoder(w).Encode(res)
}

// Run check on all groups.
func GroupsCheckAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var res Response
	groups, err := allGroups()
	if err != nil {
		res = Response{Status{500, err.Error()}, nil}
		json.NewEncoder(w).Encode(res)
		return
	}
	go check(groups)
	res = Response{Status{200, ""}, nil}
	json.NewEncoder(w).Encode(res)
}

func GroupsNotifications(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var res Response
	group, err := findGroup(mux.Vars(r)["groupId"])
	if err != nil {
		res = Response{Status{400, err.Error()}, nil}
		json.NewEncoder(w).Encode(res)
		return
	}
	notifications, err := group.allNotifications()
	if err != nil {
		res = Response{Status{500, err.Error()}, nil}
		json.NewEncoder(w).Encode(res)
		return
	}
	res = Response{Status{200, ""}, notifications}
	json.NewEncoder(w).Encode(res)

}

func GroupsShow(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var res Response
	group, err := findGroup(mux.Vars(r)["groupId"])
	if err != nil {
		res = Response{Status{400, err.Error()}, nil}
	} else {
		res = Response{Status{200, ""}, group}
	}
	json.NewEncoder(w).Encode(res)
}

func GroupConstraints(w http.ResponseWriter, r *http.Request) {
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
		res = Response{Status{500, err.Error()}, nil}
		json.NewEncoder(w).Encode(res)
		return
	}
	res = Response{Status{200, ""}, constraints}
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

func GroupsUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var res Response

	if r.Body == nil {
		res = Response{Status{400, "No body"}, nil}
		json.NewEncoder(w).Encode(res)
		return
	}

	var group Group
	err := json.NewDecoder(r.Body).Decode(&group)

	updatedGroup, err := updateGroup(mux.Vars(r)["groupId"], group)
	if err != nil {
		res = Response{Status{400, err.Error()}, nil}
		json.NewEncoder(w).Encode(res)
		return
	}

	res = Response{Status{200, ""}, updatedGroup}
	json.NewEncoder(w).Encode(res)
}
