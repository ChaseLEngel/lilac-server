package main

import (
	"fmt"
	"strconv"
)

type Request struct {
	ID              uint             `gorm:"index" gorm:"AUTO_INCREMENT" json:"request_id"`
	GroupID         uint             `json:"group_id"`
	Regex           string           `json:"regex"`
	Name            string           `json:"name"`
	DownloadPath    string           `json:"download_path"`
	RequestMachines []RequestMachine `json:"-"`
	History         []MatchHistory   `json:"-"`
}

func (group Group) insertRequest(request *Request) error {
	result := Db.Model(&group).Association("Requests").Append(request)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (group Group) allRequests() ([]Request, error) {
	var requests []Request
	result := Db.Model(&group).Related(&requests)
	if result.Error != nil {
		return nil, result.Error
	}
	return requests, nil
}

func (group Group) findRequest(id string) (Request, error) {
	requests, err := group.allRequests()
	if err != nil {
		return Request{}, err
	}
	uid, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return Request{}, err
	}

	for _, request := range requests {
		if request.ID == uint(uid) {
			return request, nil
		}
	}
	return Request{}, fmt.Errorf("record not found")
}

func (group Group) deleteRequest(id string) error {
	request, err := group.findRequest(id)
	if err != nil {
		return err
	}
	result := Db.Model(&group).Association("Requests").Delete(request)
	if result.Error != nil {
		return result.Error
	}
	Db.Delete(&request)
	return nil
}

func (group Group) updateRequest(id string, newRequest Request) (Request, error) {
	request, err := group.findRequest(id)
	if err != nil {
		return Request{}, err
	}

	result := Db.Model(&request).Updates(newRequest)
	if result.Error != nil {
		return Request{}, result.Error
	}

	return request, nil
}

func (request Request) AllRequestMachines() ([]RequestMachine, error) {
	var rms []RequestMachine
	result := Db.Model(&request).Related(&rms)
	if result.Error != nil {
		return rms, result.Error
	}
	return rms, nil
}
