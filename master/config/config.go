package config

import (
	"github.com/spf13/viper"
)

var (
	DB_HOST     string
	DB_PORT     string
	DB_USERNAME string
	DB_PASSWORD string
	DB_NAME     string
)

func init() {
	viper.SetDefault("DB_HOST", "localhost")
	DB_HOST = viper.GetString("DB_HOST")
	viper.SetDefault("DB_PORT", "3306")
	DB_PORT = viper.GetString("DB_PORT")
	viper.SetDefault("DB_USERNAME", "dev")
	DB_USERNAME = viper.GetString("DB_USERNAME")
	viper.SetDefault("DB_PASSWORD", "dev")
	DB_PASSWORD = viper.GetString("DB_PASSWORD")
	viper.SetDefault("DB_NAME", "dev")
	DB_NAME = viper.GetString("DB_NAME")
}
