package main

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strconv"
	"testing"
)

var mockGroup = Group{
	Name:          "TestGroup",
	DownloadPath:  "downloads",
	Link:          "examplefeed.com",
	Requests:      nil,
	Constraints:   nil,
	Notifications: nil,
}

func TestGroupCreateEmpty(t *testing.T) {
	setup()
	defer teardown()
	var expected = Response{
		Status: Status{400, "No body"},
		Data:   nil,
	}
	route, err := getRoute("GroupsCreate")
	if err != nil {
		t.Fatal(err)
	}
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

func TestGroupCreate(t *testing.T) {
	setup()
	defer teardown()
	var expected = Response{
		Status: Status{200, ""},
		Data:   nil,
	}
	var actual Response
	b, err := json.Marshal(mockGroup)
	if err != nil {
		t.Fatal(err)
	}
	body := bytes.NewReader(b)
	route, err := getRoute("GroupsCreate")
	if err != nil {
		t.Fatal(err)
	}
	recorder, err := setupRecorder(route, body)
	if err != nil {
		t.Fatal(err)
	}
	err = json.NewDecoder(recorder.Body).Decode(&actual)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("\nExpected %v\nActual %v", expected, actual)
	}
}

func TestGroups(t *testing.T) {
	setup()
	defer teardown()
	groups, err := allGroups()
	if err != nil {
		t.Fatal(err)
	}
	var expected = Response{
		Status: Status{200, ""},
		Data:   groups,
	}
	var actual Response
	actual.Data = new([]Group)

	route, err := getRoute("Groups")
	if err != nil {
		t.Fatal(err)
	}
	recorder, err := setupRecorder(route, nil)
	if err != nil {
		t.Fatal(err)
	}
	err = json.NewDecoder(recorder.Body).Decode(&actual)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("\nExpected %v\nActual %v", expected, actual)
	}
}

func TestGroupDelete(t *testing.T) {
	setup()
	defer teardown()
	Db.Create(&mockGroup)
	groupId := strconv.Itoa(int(mockGroup.ID))
	var expected = Response{
		Status: Status{200, ""},
		Data:   nil,
	}
	var actual Response
	route, err := getRoute("GroupsDelete")
	if err != nil {
		t.Fatal(err)
	}
	route.Path = replace(groupId, "", "", "", route.Path)
	recorder, err := setupRecorder(route, nil)
	if err != nil {
		t.Fatal(err)
	}
	err = json.NewDecoder(recorder.Body).Decode(&actual)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("\nPath %v\nExpected %v\nActual %v", route.Path, expected, actual)
	}
}

func TestGroupShow(t *testing.T) {
	setup()
	defer teardown()
	Db.Create(&mockGroup)
	groupId := strconv.Itoa(int(mockGroup.ID))
	var expected = Response{
		Status: Status{200, ""},
		Data:   &mockGroup,
	}
	var actual Response
	actual.Data = new(Group)
	route, err := getRoute("GroupsShow")
	if err != nil {
		t.Fatal(err)
	}
	route.Path = replace(groupId, "", "", "", route.Path)
	recorder, err := setupRecorder(route, nil)
	if err != nil {
		t.Fatal(err)
	}
	err = json.NewDecoder(recorder.Body).Decode(&actual)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("\nPath %v\nExpected %v\nActual %v", route.Path, expected, actual)
	}
}
