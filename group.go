package main

import "fmt"

type Group struct {
	ID            uint   `json:"id" gorm:"primary_key"`
	Name          string `json:"name"`
	DownloadPath  string `json:"downloadpath"`
	Link          string `json:"link"`
	Request       []Request
	Constraints   []Constraint
	Notifications []Notification
}

func findGroup(id string) (*Group, error) {
	group := Group{}
	result := Db.Where("ID = ?", id).Find(&group)
	if result.Error != nil {
		return nil, result.Error
	}
	return &group, nil
}

func deleteGroup(id string) error {
	result := Db.Delete(Group{}, "ID LIKE ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("record not found")
	}
	// There is a chance Gorm will delete everything so thats bad.
	if result.RowsAffected > 1 {
		return fmt.Errorf("bluk deletion")
	}
	return nil
}
