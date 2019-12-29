package config

import (
	"github.com/spf13/viper"
)

var (
	API_URL          string
	GS_DIR           string
	STORAGE_DIR      string
	TRASH_DIR        string
	HOST_TOKEN       string
	API_TOKEN        string
	GIT_COMMIT       string
	REDIS_HOST       string
	REDIS_PUB_SUB_CH string
)

func init() {
	viper.AutomaticEnv()

	viper.SetDefault("API_URL", "http://localhost:9000")
	API_URL = viper.GetString("API_URL")

	viper.SetDefault("GS_DIR", "/opt/aimpanel/gs/")
	GS_DIR = viper.GetString("GS_DIR")

	viper.SetDefault("STORAGE_DIR", "/opt/aimpanel/storage/")
	STORAGE_DIR = viper.GetString("STORAGE_DIR")

	viper.SetDefault("TRASH_DIR", "/opt/aimpanel/trash/")
	TRASH_DIR = viper.GetString("TRASH_DIR")

	viper.SetDefault("HOST_TOKEN", "")
	HOST_TOKEN = viper.GetString("HOST_TOKEN")

	viper.SetDefault("API_TOKEN", "")
	API_TOKEN = viper.GetString("API_TOKEN")

	viper.SetDefault("REDIS_HOST", "/opt/aimpanel/redis/redis.sock")
	REDIS_HOST = viper.GetString("REDIS_HOST")

	viper.SetDefault("REDIS_PUB_SUB_CH", "aimpanel")
	REDIS_PUB_SUB_CH = viper.GetString("REDIS_PUB_SUB_CH")
}
