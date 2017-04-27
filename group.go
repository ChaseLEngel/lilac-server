package main

import (
	"fmt"
	"strconv"
	"time"
)

type Group struct {
	ID            uint           `json:"group_id"`
	Name          string         `json:"name"`
	DownloadPath  string         `json:"download_path"`
	LastChecked   time.Time      `json:"last_checked"`
	Link          string         `json:"link"`
	Settings      GroupSettings  `json:"-"`
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

func (group *Group) insert() error {
	if result := Db.Create(&group); result.Error != nil {
		return result.Error
	}
	// Create cron job
	if err := master.AddSlave(int(group.ID), 30, func() { check(group.ID) }); err != nil {
		return err
	}
	// All groups have settings by default
	settings := NewGroupSettings(30, false)
	if err := group.insertGroupSettings(settings); err != nil {
		return err
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

func (group *Group) updateLastChecked() error {
	updatedGroup := group
	updatedGroup.LastChecked = time.Now()
	result := Db.Model(&group).Updates(updatedGroup)
	if result.Error != nil {
		return result.Error
	}
	return nil
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
	// Remove cron job
	intId, err := strconv.Atoi(id)
	if err != nil {
		fmt.Errorf("Failed to convert id to int")
	}
	err = master.RemoveSlave(intId)
	if err != nil {
		return fmt.Errorf("Failed to remove slave: %v", err)
	}
	return nil
}

func (group *Group) update(updated Group) error {
	result := Db.Model(&group).Updates(updated)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
