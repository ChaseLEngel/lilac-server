package main

import (
	"fmt"
	"time"
)

type Group struct {
	ID            uint           `json:"group_id"`
	Name          string         `json:"name"`
	DownloadPath  string         `json:"download_path"`
	Link          string         `json:"link"`
	Machines      []Machine      `json:"-"`
	Requests      []Request      `json:"-"`
	Constraints   []Constraint   `json:"-"`
	Notifications []Notification `json:"-"`
}

type Notification struct {
	ID        uint      `json:"notification_id" gorm:"primary_key"`
	GroupID   uint      `json:"-"`
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
}

func (group Group) allNotifications() ([]Notification, error) {
	var notifications []Notification
	result := Db.Model(&group).Related(&notifications)
	if result.Error != nil {
		return nil, result.Error
	}
	return notifications, nil
}

func insertGroup(group *Group) error {
	result := Db.Create(&group)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func allGroups() ([]Group, error) {
	groups := []Group{}
	result := Db.Find(&groups)
	if result.Error != nil {
		return nil, result.Error
	}

	return groups, nil
}

func findGroup(id string) (Group, error) {
	group := Group{}
	result := Db.Where("ID = ?", id).Find(&group)
	if result.Error != nil {
		return Group{}, result.Error
	}
	return group, nil
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

func updateGroup(id string, newGroup Group) (Group, error) {
	group, err := findGroup(id)
	if err != nil {
		return Group{}, err
	}

	result := Db.Model(&group).Updates(newGroup)
	if result.Error != nil {
		return Group{}, result.Error
	}

	return group, nil
}
