package config

import (
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

var (
	DB_HOST     string
	DB_PORT     string
	DB_USERNAME string
	DB_PASSWORD string
	DB_NAME     string

	REDIS_HOST     string
	REDIS_PASSWORD string

	GIT_COMMIT string

	DEV_MODE bool

	UPDATE_TOKEN string

	NODE_ID int64

	HTTP_DOCS_DIR     string
	HTTP_FRONTEND_DIR string
)

func init() {
	viper.AutomaticEnv()

	viper.SetDefault("DB_HOST", "localhost")
	DB_HOST = viper.GetString("DB_HOST")

	viper.SetDefault("DB_PORT", "27017")
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

	viper.SetDefault("DEV_MODE", false)
	DEV_MODE = viper.GetBool("DEV_MODE")

	viper.SetDefault("UPDATE_TOKEN", "")
	UPDATE_TOKEN = viper.GetString("UPDATE_TOKEN")

	viper.SetDefault("NODE_ID", 0)
	NODE_ID = viper.GetInt64("NODE_ID")

	binaryDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	viper.SetDefault("HTTP_DOCS_DIR", binaryDir+"/docs/")
	HTTP_DOCS_DIR = viper.GetString("HTTP_DOCS_DIR")
	viper.SetDefault("HTTP_FRONTEND_DIR", binaryDir+"/frontend/")
	HTTP_FRONTEND_DIR = viper.GetString("HTTP_FRONTEND_DIR")
}
