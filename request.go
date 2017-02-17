package main

import (
	"fmt"
	"strconv"
	"time"
)

type Request struct {
	ID           uint           `gorm:"index" gorm:"AUTO_INCREMENT" json:"request_id"`
	GroupID      uint           `json:"-"`
	Regex        string         `json:"regex"`
	Name         string         `json:"name"`
	DownloadPath string         `json:"download_path"`
	MatchCount   int            `json:"match_count"`
	Machines     []Machine      `json:"-"`
	History      []MatchHistory `json:"-"`
}

type MatchHistory struct {
	ID        uint      `gorm:"index" gorm:"AUTO_INCREMENT" json:"match_history_id"`
	RequestID uint      `json:"-"`
	Timestamp time.Time `json:"timestamp"`
	Regex     string    `json:"regex"`
	File      string    `json:"file"`
}

func (group Group) insertRequest(request Request) error {
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

func (request Request) history() (*[]MatchHistory, error) {
	var matchHistory []MatchHistory
	result := Db.Model(&request).Related(&matchHistory)
	if result.Error != nil {
		return nil, result.Error
	}
	return &matchHistory, nil
}
