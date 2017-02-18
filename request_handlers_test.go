package main

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

var mockRequest = Request{
	Name:         "TestRequest",
	DownloadPath: "/path/to/dir",
	Regex:        "TestRequest S{0-9}E{09}",
	MatchCount:   0,
	Machines:     nil,
	History:      nil,
}

func TestRequestCreateEmpty(t *testing.T) {
	setup()
	defer teardown()

	var expected = Response{
		Status: Status{400, "No body"},
		Data:   nil,
	}
	route, err := getRoute("RequestsCreate")
	if err != nil {
		t.Fatal(err)
	}
	route.Path = replace("1", "", "", "", route.Path)
	recorder, err := setupRecorder(route, nil)
	if err != nil {
		t.Fatal(err)
	}
	var actual Response
	err = json.NewDecoder(recorder.Body).Decode(&actual)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("\nExpected %v\nActual %v", expected, actual)
	}
}

func TestRequestCreate(t *testing.T) {
	setup()
	defer teardown()
	// Create new group
	var group Group
	err := insertGroup(&group)
	if err != nil {
		t.Fatal(err)
	}

	// Define what response should be returned.
	var expected = Response{
		Status: Status{200, ""},
		Data:   nil,
	}

	// Get route for request create
	route, err := getRoute("RequestsCreate")
	if err != nil {
		t.Fatal(err)
	}

	// Replace groupId parameter
	route.Path = replace("1", "", "", "", route.Path)

	// Convert new request to JSON
	b, err := json.Marshal(mockRequest)
	if err != nil {
		t.Fatal(err)
	}
	body := bytes.NewReader(b)

	// Make API request to route with new request
	recorder, err := setupRecorder(route, body)
	if err != nil {
		t.Fatal(err)
	}

	// Decode return response from API.
	var actual Response
	err = json.NewDecoder(recorder.Body).Decode(&actual)
	if err != nil {
		t.Fatal(err)
	}

	// Expected response and actual should be equal.
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("\nExpected %v\nActual %v", expected, actual)
	}
}

func TestRequests(t *testing.T) {
	setup()
	defer teardown()

	err := insertGroup(&mockGroup)
	if err != nil {
		t.Fatal(err)
	}

	err = mockGroup.insertRequest(&mockRequest)
	if err != nil {
		t.Fatal(err)
	}

	var requests []Request
	result := Db.Model(&mockGroup).Related(&requests)
	if result.Error != nil {
		t.Fatal(result.Error)
	}

	expected := Response{
		Status: Status{200, ""},
		Data:   requests,
	}

	actual := Response{
		Status: Status{},
		Data:   []Request{},
	}

	route, err := getRoute("Requests")
	if err != nil {
		t.Fatal(err)
	}

	route.Path = replace("1", "", "", "", route.Path)

	recorder, err := setupRecorder(route, nil)
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(recorder.Body.Bytes(), &actual)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expected.Status, actual.Status) {
		t.Errorf("\nExpected %v\nActual %v", expected.Status, actual.Status)
	}
}

func TestRequestHistory(t *testing.T) {
	setup()
	defer teardown()

	err := insertGroup(&mockGroup)
	if err != nil {
		t.Fatal(err)
	}

	err = mockGroup.insertRequest(&mockRequest)
	if err != nil {
		t.Fatal(err)
	}

	history := MatchHistory{
		Timestamp: time.Now(),
		Regex:     "TestRequest S{0-9}E{0-9}",
		File:      "TestRequest S1E3",
	}

	err = mockRequest.insertMatchHistory(&history)

	matchHistorys, err := mockRequest.history()
	if err != nil {
		t.Fatal(err)
	}

	expected := Response{
		Status: Status{200, ""},
		Data:   matchHistorys,
	}

	route, err := getRoute("RequestsHistory")
	if err != nil {
		t.Fatal(err)
	}
	route.Path = replace("1", "1", "", "", route.Path)

	recorder, err := setupRecorder(route, nil)
	if err != nil {
		t.Fatal(err)
	}

	var actual Response
	actual.Data = []MatchHistory{}
	err = json.Unmarshal(recorder.Body.Bytes(), &actual)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expected.Status, actual.Status) {
		t.Errorf("\nExpected %#v\nActual %#v", expected.Status, actual.Status)
	}
}
