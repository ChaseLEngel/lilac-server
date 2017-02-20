package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
)

// Sets up http recorder for specific Route.
func setupRecorder(route Route, body io.Reader) (*httptest.ResponseRecorder, error) {
	req, err := http.NewRequest(route.Method, route.Path, body)
	if err != nil {
		return nil, err
	}
	recorder := httptest.NewRecorder()
	// Need to call Gorrila router setup or else http test wont work.
	NewRouter().ServeHTTP(recorder, req)
	return recorder, nil
}

func setup() {
	var err error
	err = initDatabase("testing.db")
	if err != nil {
		panic("Failed to init the database.")
	}
}

func teardown() {
	os.Remove("testing.db")
}

// Replace endpoint parameters with values
func replace(groupId, requestId, machineId, file, endpoint string) string {
	paras := map[string]string{
		"{groupId}":   groupId,
		"{requestId}": requestId,
		"{machineId}": machineId,
		"{file}":      file,
	}
	for key, val := range paras {
		endpoint = strings.Replace(endpoint, key, val, -1)
	}
	return endpoint
}

// In order to handle Response's Data interface we need test with a json.RawMessage
type ResponseTesting struct {
	Status Status          `json:"status"`
	Data   json.RawMessage `json:"data"`
}
