package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var Db *gorm.DB

func initDatabase() error {
	var err error
	Db, err = gorm.Open("sqlite3", "lilac.db")
	if err != nil {
		return err
	}

	// Update or create tables
	Db.AutoMigrate(
		&Constraint{},
		&Group{},
		&Machine{},
		&Notification{},
		&Request{},
		&MatchHistory{},
	)
	return nil
}
