package main

import (
	"fmt"
	"strconv"
	"time"
)

type MatchHistory struct {
	ID        uint      `gorm:"index" gorm:"AUTO_INCREMENT" json:"match_history_id"`
	RequestID uint      `json:"request_id"`
	Timestamp time.Time `json:"timestamp"`
	Title     string    `json:"title"` // RSS item title
	Name      string    `json:"name"`  // Name of torrent file
	Files     string    `json:"files"` // Comma seperated list of files torrent contains
	Size      int       `json:"size"`  // Total size of torrent files in bytes
}

func (request Request) insertMatchHistory(history *MatchHistory) error {
	result := Db.Model(&request).Association("History").Append(history)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (request Request) deleteMatchHistory(id string) (MatchHistory, error) {
	history, err := request.findMatchHistory(id)
	if err != nil {
		return MatchHistory{}, err
	}
	result := Db.Model(&request).Association("History").Delete(history)
	if result.Error != nil {
		return MatchHistory{}, result.Error
	}
	Db.Delete(&history)
	return history, nil
}

func (request Request) history() ([]MatchHistory, error) {
	var matchHistory []MatchHistory
	result := Db.Model(&request).Related(&matchHistory)
	if result.Error != nil {
		return nil, result.Error
	}
	return matchHistory, nil
}

func (request Request) findMatchHistory(id string) (MatchHistory, error) {
	history, err := request.history()
	if err != nil {
		return MatchHistory{}, err
	}
	uid, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return MatchHistory{}, err
	}
	for _, h := range history {
		if h.ID == uint(uid) {
			return h, nil
		}
	}
	return MatchHistory{}, fmt.Errorf("record not found")
}
