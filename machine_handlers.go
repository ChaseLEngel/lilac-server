package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func Machines(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var res Response
	machines, err := allMachines()
	if err != nil {
		res = Response{Status{400, err.Error()}, nil}
		json.NewEncoder(w).Encode(res)
		return
	}
	res = Response{Status{200, ""}, machines}
	json.NewEncoder(w).Encode(res)
	return
}

func MachinesCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var res Response
	if r.Body == nil {
		res = Response{Status{400, "No body"}, nil}
		json.NewEncoder(w).Encode(res)
		return
	}

	machine := new(Machine)
	json.NewDecoder(r.Body).Decode(&machine)
	err := machine.insert()
	if err != nil {
		res = Response{Status{500, err.Error()}, nil}
		json.NewEncoder(w).Encode(res)
		return
	}
	res = Response{Status{200, ""}, machine}
	json.NewEncoder(w).Encode(res)
}

func MachinesUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var res Response

	var newMachine Machine
	err := json.NewDecoder(r.Body).Decode(&newMachine)
	if err != nil {
		res = Response{Status{400, err.Error()}, nil}
		json.NewEncoder(w).Encode(res)
		return
	}

	machine, err := findMachine(mux.Vars(r)["machineId"])
	if err != nil {
		res = Response{Status{400, err.Error()}, nil}
		json.NewEncoder(w).Encode(res)
		return
	}

	err = machine.update(newMachine)
	if err != nil {
		res = Response{Status{400, err.Error()}, nil}
		json.NewEncoder(w).Encode(res)
		return
	}

	res = Response{Status{200, ""}, machine}
	json.NewEncoder(w).Encode(res)
}

func MachinesDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var res Response

	machine, err := findMachine(mux.Vars(r)["machineId"])
	if err != nil {
		res = Response{Status{400, err.Error()}, nil}
		json.NewEncoder(w).Encode(res)
		return
	}

	machine.delete()

	res = Response{Status{200, ""}, machine}
	json.NewEncoder(w).Encode(res)
}
