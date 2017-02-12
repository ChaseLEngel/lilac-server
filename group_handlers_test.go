package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

var mockGroup = Group{
	Name:          "TestGroup",
	DownloadPath:  "downloads",
	Link:          "examplefeed.com",
	Request:       nil,
	History:       nil,
	Constraints:   nil,
	Notifications: nil,
}

func TestMain(m *testing.M) {
	var err error
	err = initDatabase()
	if err != nil {
		panic("Failed to init the database.")
	}

	m.Run()
}

func TestGroupsCreateEmpty(t *testing.T) {
	var expected = Response{
		Status: Status{400, "No body"},
		Data:   nil,
	}
	req, err := http.NewRequest("POST", "/groups", nil)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(GroupsCreate)
	handler.ServeHTTP(recorder, req)

	var actual Response
	err = json.NewDecoder(recorder.Body).Decode(&actual)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("\nExpected %v\nActual %v", expected, actual)
	}
}

func TestGroupsCreate(t *testing.T) {
	var expected = Response{
		Status: Status{200, ""},
		Data:   nil,
	}
	b, err := json.Marshal(mockGroup)
	if err != nil {
		t.Fatal(err)
	}
	body := bytes.NewReader(b)
	req, err := http.NewRequest("POST", "/groups", body)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(GroupsCreate)
	handler.ServeHTTP(recorder, req)
	var actual Response
	err = json.NewDecoder(recorder.Body).Decode(&actual)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("\nExpected %v\nActual %v", expected, actual)
	}
}

func TestGroups(t *testing.T) {
	var expected = Response{
		Status: Status{200, ""},
		Data:   &mockGroup,
	}
	req, err := http.NewRequest("GET", "/groups", nil)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(Groups)
	handler.ServeHTTP(recorder, req)
	var actual Response
	actual.Data = new(Group)
	err = json.NewDecoder(recorder.Body).Decode(&actual)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("\nExpected %v\nActual %v", expected, actual)
	}

}
