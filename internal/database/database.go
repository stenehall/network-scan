package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB model.
type DB struct {
	db *gorm.DB
}

type Host struct {
	IP       string
	Hostname string
}

// AddIfNotExist adds the ip/hostname if not already added.
func (database DB) AddIfNotExist(ipAddr string, hostname string) (tx *gorm.DB) {
	var hosts []Host

	result := database.db.Where("IP = ?", ipAddr).First(&hosts)

	if result.Error != nil {
		// Save the new IP to DB.
		database.db.Create(&Host{IP: ipAddr, Hostname: hostname})
	}

	return result
}

func (database DB) GetAll() []Host {
	var hosts []Host
	database.db.Find(&hosts)

	return hosts
}

// Database constructor.
func Database(dbName string) DB {
	var hosts []Host

	gormInstance, err := gorm.Open(
		sqlite.Open(dbName),
		&gorm.Config{
			SkipDefaultTransaction:                   false,
			NamingStrategy:                           nil,
			FullSaveAssociations:                     false,
			Logger:                                   logger.Default.LogMode(logger.Silent),
			NowFunc:                                  nil,
			DryRun:                                   false,
			PrepareStmt:                              false,
			DisableAutomaticPing:                     false,
			DisableForeignKeyConstraintWhenMigrating: false,
			DisableNestedTransaction:                 false,
			AllowGlobalUpdate:                        false,
			QueryFields:                              false,
			CreateBatchSize:                          0,
			ClauseBuilders:                           nil,
			ConnPool:                                 nil,
			Dialector:                                nil,
			Plugins:                                  nil,
		},
	)
	if err != nil {
		panic("failed to connect database")
	}

	migrationError := gormInstance.AutoMigrate(&hosts)
	if migrationError != nil {
		panic("failed to auto-migrate database")
	}

	return DB{
		gormInstance,
	}
}
