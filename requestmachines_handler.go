package main

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
)

func RequestMachines(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var res Response

	var request Request
	result := Db.Find(&request, "ID = "+mux.Vars(r)["requestID"])
	if result.Error != nil {
		res = NewResponse(400, result.Error, nil)
		json.NewEncoder(w).Encode(res)
		return
	}

	requestMachines, err := request.AllRequestMachines()
	if err != nil {
		res = NewResponse(400, err, nil)
		json.NewEncoder(w).Encode(res)
		return
	}

	res = NewResponse(200, nil, requestMachines)
	json.NewEncoder(w).Encode(res)
}

type MachineDestination struct {
	MachineID   string `json:"machine_id"`
	Destination string `json:"destination"`
}

func RequestMachinesCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var res Response

	log.Info(r.Body)

	if r.Body == nil {
		res = NewResponse(400, errors.New("No Body"), nil)
		json.NewEncoder(w).Encode(res)
		return
	}

	var request Request
	result := Db.Find(&request, "ID = "+mux.Vars(r)["requestID"])
	if result.Error != nil {
		res = NewResponse(400, result.Error, nil)
		json.NewEncoder(w).Encode(res)
		return
	}

	Db.Model(&request).Association("RequestMachines").Clear()

	var machineDestinations []MachineDestination
	err := json.NewDecoder(r.Body).Decode(&machineDestinations)
	if err != nil {
		res = NewResponse(400, err, nil)
		json.NewEncoder(w).Encode(res)
		return
	}

	for _, md := range machineDestinations {
		machine, err := findMachine(md.MachineID)
		if err != nil {
			res = NewResponse(400, err, nil)
			json.NewEncoder(w).Encode(res)
			return
		}
		var rm RequestMachine
		rm.MachineID = machine.ID
		rm.Destination = md.Destination
		Db.Model(&request).Association("RequestMachines").Append(&rm)
	}

	var requestMachines []RequestMachine
	result = Db.Model(&request).Related(&requestMachines)
	if result.Error != nil {
		res = NewResponse(400, result.Error, nil)
		json.NewEncoder(w).Encode(res)
		return
	}

	res = NewResponse(200, nil, requestMachines)
	json.NewEncoder(w).Encode(res)
}
