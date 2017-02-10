package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Replace endpoint parameters with values
func replace(endpoint string) string {
	apiVars := map[string]string{
		"{groupId}":   "1",
		"{requestId}": "1",
		"{machineId}": "1",
		"{file}":      "testfile.txt",
	}
	for key, val := range apiVars {
		newStr := strings.Replace(endpoint, key, val, -1)
		if newStr != endpoint {
			return newStr
		}
	}
	return endpoint
}

// Test that all routes return application/json in Header
func TestContentTypeIsJSON(t *testing.T) {
	for _, route := range routes {
		endpoint := replace(route.Path)
		req, err := http.NewRequest(route.Method, endpoint, nil)
		if err != nil {
			t.Fatal(err)
		}
		recorder := httptest.NewRecorder()
		handler := http.HandlerFunc(route.HandlerFunc)
		handler.ServeHTTP(recorder, req)
		if contentType := recorder.Header().Get("Content-Type"); contentType != "application/json" {
			t.Errorf("Content-Type for %v was %v not application/json", endpoint, contentType)
		}
	}
}
