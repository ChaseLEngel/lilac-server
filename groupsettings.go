package main

type GroupSettings struct {
	ID           uint `json:"settings_id" gorm:"primary_key"`
	GroupId      uint `json:"group_id"`
	Interval     int  `json:"interval"`      // RSS checking interval
	AutoTransfer bool `json:"auto_transfer"` // Should files be transfered after download
}

func NewGroupSettings(interval int, autoTransfer bool) GroupSettings {
	settings := new(GroupSettings)
	settings.Interval = interval
	settings.AutoTransfer = autoTransfer
	return *settings
}

func (group Group) insertGroupSettings(settings GroupSettings) error {
	result := Db.Model(&group).Association("Settings").Replace(settings)
	if result.Error != nil {
		return result.Error
	}
	// Assume the interval has changed and
	// update the cron job with new interval.
	err := master.ChangeTime(int(group.ID), settings.Interval)
	if err != nil {
		return err
	}
	return nil
}

func (group Group) GroupSettings() (GroupSettings, error) {
	var settings GroupSettings
	result := Db.Model(&group).Related(&settings)
	if result.Error != nil {
		return GroupSettings{}, result.Error
	}
	return settings, nil
}

func (settings GroupSettings) update(interval int, autoTransfer bool) error {
	result := Db.Model(&settings).Updates(GroupSettings{Interval: interval, AutoTransfer: autoTransfer})
	if result.Error != nil {
		return result.Error
	}
	// Assume the interval has changed and
	// update the cron job with new interval.
	err := master.ChangeTime(int(settings.GroupId), interval)
	if err != nil {
		return err
	}
	return nil
}
