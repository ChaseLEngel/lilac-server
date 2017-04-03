package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func RequestMachines(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var res Response

	var request Request
	result := Db.Find(&request, "ID = "+mux.Vars(r)["requestID"])
	if result.Error != nil {
		res = Response{Status{400, result.Error.Error()}, nil}
		json.NewEncoder(w).Encode(res)
		return
	}

	requestMachines, err := request.AllRequestMachines()
	if err != nil {
		res = Response{Status{400, err.Error()}, nil}
		json.NewEncoder(w).Encode(res)
		return
	}

	res = Response{Status{200, ""}, requestMachines}
	json.NewEncoder(w).Encode(res)
}

type MachineDestination struct {
	MachineID   string `json:"machine_id"`
	Destination string `json:"destination"`
}

func RequestMachinesCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var res Response

	if r.Body == nil {
		res = Response{Status{400, "No body"}, nil}
		json.NewEncoder(w).Encode(res)
		return
	}

	var request Request
	result := Db.Find(&request, "ID = "+mux.Vars(r)["requestID"])
	if result.Error != nil {
		res = Response{Status{400, result.Error.Error()}, nil}
		json.NewEncoder(w).Encode(res)
		return
	}

	Db.Model(&request).Association("RequestMachines").Clear()

	var machineDestinations []MachineDestination
	err := json.NewDecoder(r.Body).Decode(&machineDestinations)
	if err != nil {
		res = Response{Status{400, err.Error()}, nil}
		json.NewEncoder(w).Encode(res)
		return
	}

	for _, md := range machineDestinations {
		machine, err := findMachine(md.MachineID)
		if err != nil {
			res = Response{Status{400, err.Error()}, nil}
			json.NewEncoder(w).Encode(res)
			return
		}
		destination := md.Destination
		var rm RequestMachine
		rm.Machine = *machine
		rm.Destination = destination
		Db.Model(&request).Association("RequestMachines").Append(&rm)
	}

	var requestMachines []RequestMachine
	result = Db.Model(&request).Related(&requestMachines)
	if result.Error != nil {
		res = Response{Status{400, result.Error.Error()}, nil}
		json.NewEncoder(w).Encode(res)
		return
	}

	res = Response{Status{200, ""}, requestMachines}
	json.NewEncoder(w).Encode(res)
}

func RequestMachinesHistory(w http.ResponseWriter, r *http.Request) {
}
