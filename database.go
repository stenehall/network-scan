package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DB struct {
	db *gorm.DB
}

// @TODO Were do we want this? Outside the function?
type Host struct {
	gorm.Model
	IP       string
	Hostname string
}

// @TODO Were do we want this? Capitalized?
var Hosts []Host

func (database DB) Check(ip string, hostname string) (tx *gorm.DB) {
	result := database.db.Where("IP = ?", ip).First(&Hosts)

	if result.Error != nil {
		// Save the new IP to DB.
		database.db.Create(&Host{IP: ip, Hostname: hostname})
	}

	return result
}

func Database(db_name string) DB {
	db, err := gorm.Open(sqlite.Open(db_name), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Hosts)

	return DB{
		db,
	}
}
