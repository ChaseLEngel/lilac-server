package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"regexp"
)

type Auth struct {
	Token string `json:"token"`
}

type User struct {
	User     string
	Password string
}

func login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var res Response
	if r.Body == nil {
		res = NewResponse(400, fmt.Errorf("No body"), nil)
	}

	var user User
	json.NewDecoder(r.Body).Decode(&user)
	if conf.User == user.User && conf.Password == user.Password {
		res = NewResponse(200, nil, Auth{Token: jwtData.TokenString})
		log.Infof("Successful login for user: \"%v\"\n", user.User)
	} else {
		res = NewResponse(401, fmt.Errorf("Bad credentials for user: \"%v\"", user.User), nil)
	}
	json.NewEncoder(w).Encode(res)
}

type File struct {
	File string
}

func transferRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var res Response
	groups, err := allGroups()
	if err != nil {
		res = Response{Status{400, err.Error()}, nil}
	}

	if r.Body == nil {
		res = Response{Status{400, "No body"}, nil}
	}

	var file File
	json.NewDecoder(r.Body).Decode(&file)
	for _, group := range groups {
		requests, err := group.allRequests()
		if err != nil {
			log.Errorf("Failed to get allRequests for group %v err: %v\n", group.ID, err)
			continue
		}
		for _, request := range requests {
			log.Infof("Comparing %v with %v", filepath.Base(file.File), request.Regex)
			if matched, err := regexp.MatchString(request.Regex, filepath.Base(file.File)); !matched || err != nil {
				continue
			}
			log.Infof("Sending %v for %v", file.File, request.Name)
			err := send(request, file.File)
			if err != nil {
				log.Errorf("Failed to send %v err: %v\n", file.File, err)
				res = Response{Status{500, err.Error()}, nil}
			}
		}
	}
	if err == nil {
		res = Response{Status{200, ""}, nil}
	}
	json.NewEncoder(w).Encode(res)
}
