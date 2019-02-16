package config

import (
	"github.com/spf13/viper"
)

var (
	RABBITMQ_HOST     string
	RABBITMQ_PORT     string
	RABBITMQ_USERNAME string
	RABBITMQ_PASSWORD string
)

func init() {
	viper.SetDefault("RABBITMQ_HOST", "localhost")
	RABBITMQ_HOST = viper.GetString("RABBITMQ_HOST")
	viper.SetDefault("RABBITMQ_PORT", "5672")
	RABBITMQ_PORT = viper.GetString("RABBITMQ_PORT")
	viper.SetDefault("RABBITMQ_USERNAME", "guest")
	RABBITMQ_USERNAME = viper.GetString("RABBITMQ_USERNAME")
	viper.SetDefault("RABBITMQ_PASSWORD", "guest")
	RABBITMQ_PASSWORD = viper.GetString("RABBITMQ_PASSWORD")
}
