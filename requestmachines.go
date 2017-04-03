package main

import (
	"fmt"
	"time"
)

type RequestMachine struct {
	ID          uint   `json:"request_machine_id" gorm:"primary_key`
	RequestID   uint   `json:"request_id"`
	MachineID   uint   `json:"machine_id"`
	Destination string `json:"destination"`
}

type TransferHistory struct {
	ID               uint `json:"transfer_history_id" gorm:"primary_key"`
	RequestMachineId uint `json:"request_machine_id"`
	Timestamp        time.Time
	Destination      string
	Filename         string
}

func (rm *RequestMachine) insert(machine Machine, destination string) error {
	rm.Destination = destination

	result := Db.Create(rm)
	if result.Error != nil {
		return result.Error
	}

	Db.Model(rm).Association("Machine").Replace(machine)
	return nil
}

func (rm *RequestMachine) delete() error {
	result := Db.Delete(rm)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("record not found")
	}
	return nil
}

func (rm *RequestMachine) update(updated RequestMachine) error {
	result := Db.Model(rm).Updates(updated)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
