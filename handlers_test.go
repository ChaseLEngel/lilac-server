package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// Test that all routes return application/json in Header
func TestContentTypeIsJSON(t *testing.T) {
	for _, route := range routes {
		endpoint := replace("", "", "", "", route.Path)
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
