package config

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	RABBITMQ_HOST     string
	RABBITMQ_PORT     string
	RABBITMQ_USERNAME string
	RABBITMQ_PASSWORD string
	GS_DIR            string
	STORAGE_DIR       string
)

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		logrus.Error(fmt.Sprintf("Fatal error config file: %s", err))
	}

	viper.SetDefault("RABBITMQ_HOST", "localhost")
	RABBITMQ_HOST = viper.GetString("RABBITMQ_HOST")

	viper.SetDefault("RABBITMQ_PORT", "5672")
	RABBITMQ_PORT = viper.GetString("RABBITMQ_PORT")

	viper.SetDefault("RABBITMQ_USERNAME", "guest")
	RABBITMQ_USERNAME = viper.GetString("RABBITMQ_USERNAME")

	viper.SetDefault("RABBITMQ_PASSWORD", "guest")
	RABBITMQ_PASSWORD = viper.GetString("RABBITMQ_PASSWORD")

	viper.SetDefault("GS_DIR", "/opt/aimpanel/gs/")
	GS_DIR = viper.GetString("GS_DIR")

	viper.SetDefault("STORAGE_DIR", "/opt/aimpanel/storage/")
	STORAGE_DIR = viper.GetString("STORAGE_DIR")
}
