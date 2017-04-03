package main

import (
	"fmt"
)

type Machine struct {
	ID   uint   `json:"machine_id" gorm:"primary_key"`
	Host string `json:"host"`
	Port string `json:"port"`
	User string `json:"user"`
}

func allMachines() ([]Machine, error) {
	var machines []Machine
	result := Db.Find(&machines)
	if result.Error != nil {
		return nil, result.Error
	}
	return machines, nil
}

func findMachine(id string) (*Machine, error) {
	machine := new(Machine)
	result := Db.Where("ID = ?", id).Find(machine)
	if result.Error != nil {
		return nil, result.Error
	}
	return machine, nil

}

func (machine *Machine) insert() error {
	result := Db.Create(&machine)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (machine *Machine) delete() error {
	result := Db.Delete(machine)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("record not found")
	}
	return nil
}

func (machine *Machine) update(newMachine Machine) error {
	result := Db.Model(&machine).Updates(newMachine)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
