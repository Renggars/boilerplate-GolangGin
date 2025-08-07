package config

import (
	"restApi-GoGin/src/models"

	"gorm.io/gorm"
)

func RunMigration(db *gorm.DB) {
	db.AutoMigrate(
		&models.User{},
	)
}
