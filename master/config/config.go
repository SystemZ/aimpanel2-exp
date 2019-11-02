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

	REDIS_HOST     string
	REDIS_PASSWORD string

	RABBITMQ_TLS           bool
	RABBITMQ_HOST          string
	RABBITMQ_HOST_EXTERNAL string
	RABBITMQ_PORT          string
	RABBITMQ_PORT_EXTERNAL string
	RABBITMQ_PORT_API      string
	RABBITMQ_USERNAME      string
	RABBITMQ_PASSWORD      string
	RABBITMQ_VHOST         string

	LOG_DEV bool
)

func init() {
	viper.AutomaticEnv()

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

	viper.SetDefault("REDIS_HOST", "localhost")
	REDIS_HOST = viper.GetString("REDIS_HOST")

	viper.SetDefault("REDIS_PASSWORD", "")
	REDIS_PASSWORD = viper.GetString("REDIS_PASSWORD")

	viper.SetDefault("RABBITMQ_TLS", false)
	RABBITMQ_TLS = viper.GetBool("RABBITMQ_TLS")

	viper.SetDefault("RABBITMQ_HOST", "localhost")
	RABBITMQ_HOST = viper.GetString("RABBITMQ_HOST")

	viper.SetDefault("RABBITMQ_HOST_EXTERNAL", "localhost")
	RABBITMQ_HOST_EXTERNAL = viper.GetString("RABBITMQ_HOST_EXTERNAL")

	viper.SetDefault("RABBITMQ_PORT", "5672")
	RABBITMQ_PORT = viper.GetString("RABBITMQ_PORT")

	viper.SetDefault("RABBITMQ_PORT_EXTERNAL", "5672")
	RABBITMQ_PORT_EXTERNAL = viper.GetString("RABBITMQ_PORT_EXTERNAL")

	viper.SetDefault("RABBITMQ_PORT_API", "15672")
	RABBITMQ_PORT_API = viper.GetString("RABBITMQ_PORT_API")

	viper.SetDefault("RABBITMQ_USERNAME", "guest")
	RABBITMQ_USERNAME = viper.GetString("RABBITMQ_USERNAME")

	viper.SetDefault("RABBITMQ_PASSWORD", "guest")
	RABBITMQ_PASSWORD = viper.GetString("RABBITMQ_PASSWORD")

	viper.SetDefault("RABBITMQ_VHOST", "/")
	RABBITMQ_VHOST = viper.GetString("RABBITMQ_VHOST")

	viper.SetDefault("LOG_DEV", false)
	LOG_DEV = viper.GetBool("LOG_DEV")
}
