package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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
