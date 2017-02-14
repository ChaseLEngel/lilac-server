package main

import (
	"fmt"
)

type Group struct {
	ID            uint           `json:"group_id"`
	Name          string         `json:"name"`
	DownloadPath  string         `json:"download_path"`
	Link          string         `json:"link"`
	Requests      []Request      `json:"-"`
	Constraints   []Constraint   `json:"-"`
	Notifications []Notification `json:"-"`
}

func insertGroup() error {
	return nil

}

func allGroups() (*[]Group, error) {
	groups := []Group{}
	result := Db.Find(&groups)
	if result.Error != nil {
		return nil, result.Error
	}

	return &groups, nil
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
	return nil
}
