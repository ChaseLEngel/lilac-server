package main

import (
	"fmt"
	"strconv"
)

type Machine struct {
	ID          uint   `json:"machine_id" gorm:"primary_key"`
	GroupID     uint   `json:"-"`
	Hostname    string `json:"hostname"`
	Credentials string `json:"credentials"`
}

func (group Group) allMachines() ([]Machine, error) {
	var machines []Machine
	result := Db.Model(&group).Related(&machines)
	if result.Error != nil {
		return nil, result.Error
	}
	return machines, nil
}

func (group Group) insertMachine(machine *Machine) error {
	result := Db.Model(&group).Association("Machines").Append(machine)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (group Group) findMachine(id string) (Machine, error) {
	machines, err := group.allMachines()
	if err != nil {
		return Machine{}, err
	}

	uid, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return Machine{}, err
	}
	for _, machine := range machines {
		if machine.ID == uint(uid) {
			return machine, nil
		}
	}
	return Machine{}, fmt.Errorf("record not found")
}

func (group Group) deleteMachine(id string) error {
	machine, err := group.findMachine(id)
	if err != nil {
		return err
	}
	result := Db.Model(&group).Association("Machines").Delete(machine)
	if result.Error != nil {
		return result.Error
	}
	Db.Delete(&machine)
	return nil
}

func (group Group) updateMachine(id string, newMachine Machine) (Machine, error) {
	machine, err := group.findMachine(id)
	if err != nil {
		return Machine{}, err
	}

	result := Db.Model(&machine).Updates(newMachine)
	if result.Error != nil {
		return Machine{}, result.Error
	}

	return machine, nil
}
