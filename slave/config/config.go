package config

import (
	"github.com/spf13/viper"
)

var (
	API_URL     string
	GS_DIR      string
	STORAGE_DIR string
	TRASH_DIR   string
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
}
