package main

import (
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
	Timestamp time.Time
	Regex     string
	File      string
}

func (group Group) insertRequest(request Request) error {
	result := Db.Model(&group).Association("Requests")
	if result.Error != nil {
		return result.Error
	}
	group.Requests = append(group.Requests, request)
	result = Db.Model(&group).Association("Requests").Append(request)
	if result.Error != nil {
		return result.Error
	}
	Db.Save(&group)
	return nil
}

func (group Group) allRequests() (*[]Request, error) {
	var requests []Request
	result := Db.Model(&group).Related(&requests)
	if result.Error != nil {
		return nil, result.Error
	}
	return &requests, nil
}
