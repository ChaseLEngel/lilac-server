package main

import (
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
	Db.DropTable(&Group{})
	os.Remove("testing.db")
	defer Db.Close()
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
		newStr := strings.Replace(endpoint, key, val, -1)
		if newStr != endpoint {
			return newStr
		}
	}
	return endpoint
}
