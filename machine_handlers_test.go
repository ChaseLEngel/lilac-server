package main

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestMachines(t *testing.T) {
	setup()
	defer teardown()
	var group Group
	err := group.insert()
	if err != nil {
		t.Fatal(err)
	}
	var machine Machine
	err = group.insertMachine(&machine)
	if err != nil {
		t.Fatal(err)
	}
	machines, err := group.allMachines()
	if err != nil {
		t.Fatal(err)
	}

	route, err := getRoute("Machines")
	if err != nil {
		t.Fatal(err)
	}

	route.Path = replace("1", "", "", "", route.Path)

	recorder, err := setupRecorder(route, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Get rid of struct values that are ignored by JSON.
	b, err := json.Marshal(NewResponse(200, nil, machines))
	if err != nil {
		t.Fatal(err)
	}

	var expected ResponseTesting
	json.Unmarshal(b, &expected)

	var expectedData []Machine
	json.Unmarshal([]byte(expected.Data), &expectedData)

	var actual ResponseTesting
	json.Unmarshal(recorder.Body.Bytes(), &actual)

	var actualData []Machine
	json.Unmarshal([]byte(actual.Data), &actualData)

	if !reflect.DeepEqual(expected.Status, actual.Status) {
		t.Errorf("\nPath: %v\nExpected: %#v\nActual: %#v\n", route.Path, expected, actual)
	}

	if !reflect.DeepEqual(expectedData, actualData) {
		t.Errorf("\nPath: %v\nExpected: %#v\nActual: %#v\n", route.Path, expected.Data, actualData)
	}
}
