package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Replace /script/checkin/{scriptName}
// with /script/checkin/test for all routes
// If endpoint doesn't have any replacments return unchanged endpoint
func replace(endpoint string) string {
	apiVars := map[string]string{
		"{scriptId}":   "1",
		"{scriptName}": "test",
		"{mailId}":     "test@example.com",
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
