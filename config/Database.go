package config

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitMySQL(config *config) *gorm.DB {
	conString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=%s",
		config.Db.DB_USERNAME,
		config.Db.DB_PASSWORD,
		config.Db.DB_HOST,
		config.Db.DB_PORT,
		config.Db.DB_DATABASE,
		config.App.APP_TIMEZONE,
	)

	db, err := gorm.Open(mysql.Open(conString), &gorm.Config{})

	if err != nil {
		log.Fatal("Error while connecting to database", err)
	}

	return db
}
