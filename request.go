package main

import (
	_ "github.com/jinzhu/gorm"
	"time"
)

type Request struct {
	RequestID    int `gorm:"index"`
	Regex        string
	Name         string
	DownloadPath string
	GroupRefer   uint
	MatchCount   int
	Machines     []Machine
	History      []MatchHistory
}

type MatchHistory struct {
	Timestamp time.Time
	Regex     string
	File      string
}

func (group Group) insertRequest() (*[]Request, error) {
	return nil, nil
}

func (group Group) allRequests() (*[]Request, error) {
	var requests []Request
	result := Db.Model(&group).Related(&requests)
	if result.Error != nil {
		return nil, result.Error
	}
	return &requests, nil
}
