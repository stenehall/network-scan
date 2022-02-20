package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DB struct {
	db *gorm.DB
}

type host struct {
	gorm.Model
	IP       string
	Hostname string
}

var hosts []host

func (database DB) AddIfNotExist(ip string, hostname string) (tx *gorm.DB) {
	result := database.db.Where("IP = ?", ip).First(&hosts)

	if result.Error != nil {
		// Save the new IP to DB.
		database.db.Create(&host{IP: ip, Hostname: hostname})
	}

	return result
}

func Database(dbName string) DB {
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic("failed to connect database")
	}

	migrationError := db.AutoMigrate(&hosts)
	if migrationError != nil {
		panic("failed to auto-migrate database")
	}

	return DB{
		db,
	}
}
