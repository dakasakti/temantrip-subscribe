package config

import (
	"fmt"
	"os"
	"sync"

	"github.com/spf13/viper"
)

type config struct {
	App struct {
		APP_NAME     string
		APP_URL      string
		APP_PORT     string
		APP_TIMEZONE string
	}
	Db struct {
		DB_CONNECTION string
		DB_HOST       string
		DB_PORT       string
		DB_DATABASE   string
		DB_USERNAME   string
		DB_PASSWORD   string
	}
	Mail struct {
		API string
	}
}

var lock = &sync.Mutex{}
var appConfig *config

func GetConfig() *config {
	lock.Lock()
	defer lock.Unlock()

	if appConfig == nil {
		appConfig = readInConfig()
	}

	return appConfig
}

func readInConfig() *config {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	config := new(config)
	config.App.APP_NAME = getEnv("APP_NAME", viper.GetString("APP_NAME"))
	config.App.APP_URL = getEnv("APP_URL", viper.GetString("APP_URL"))
	config.App.APP_PORT = getEnv("PORT", viper.GetString("APP_PORT"))
	config.App.APP_TIMEZONE = getEnv("APP_TIMEZONE", viper.GetString("APP_TIMEZONE"))

	config.Db.DB_CONNECTION = getEnv("DB_CONNECTION", viper.GetString("DB_CONNECTION"))
	config.Db.DB_HOST = getEnv("DB_HOST", viper.GetString("DB_HOST"))
	config.Db.DB_PORT = getEnv("DB_PORT", viper.GetString("DB_PORT"))
	config.Db.DB_DATABASE = getEnv("DB_DATABASE", viper.GetString("DB_DATABASE"))
	config.Db.DB_USERNAME = getEnv("DB_USERNAME", viper.GetString("DB_USERNAME"))
	config.Db.DB_PASSWORD = getEnv("DB_PASSWORD", viper.GetString("DB_PASSWORD"))

	config.Mail.API = getEnv("MAIL_API", viper.GetString("MAIL_API"))

	fmt.Println(config)
	return config
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}
