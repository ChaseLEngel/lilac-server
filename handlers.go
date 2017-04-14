package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"regexp"
)

func login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
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
			fmt.Printf("Failed to get allRequests for group %v err: %v\n", group.ID, err)
			continue
		}
		for _, request := range requests {
			if matched, err := regexp.MatchString(request.Regex, filepath.Base(file.File)); !matched || err != nil {
				continue
			}
			err := send(request, file.File)
			if err != nil {
				fmt.Println("Failed to send %v err: %v\n", file.File, err)
				res = Response{Status{500, err.Error()}, nil}
			}
		}
	}
	if err == nil {
		res = Response{Status{200, ""}, nil}
	}
	json.NewEncoder(w).Encode(res)
}
