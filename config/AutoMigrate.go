package config

import (
	"temantrip-subscribe/entities"

	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(&entities.User{})
}
