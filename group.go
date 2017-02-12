package main

import "reflect"

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
	if reflect.DeepEqual(Group{}, result.Value) {
		return nil, nil
	}
	return &group, nil
}
