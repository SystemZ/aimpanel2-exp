package config

import (
	"github.com/spf13/viper"
)

var (
	GS_DIR                      string
	STORAGE_DIR                 string
	TRASH_DIR                   string
	BACKUP_DIR                  string
	HOST_TOKEN                  string
	GIT_COMMIT                  string
	REDIS_HOST                  string
	REDIS_USERNAME              string
	REDIS_PASSWORD              string
	REDIS_PUB_SUB_AGENT_CH      string
	REDIS_PUB_SUB_WRAPPER_CH    string
	REDIS_PUB_SUB_SUPERVISOR_CH string
	MASTER_URLS                 []string
	HW_ID                       string
)

func init() {
	viper.AutomaticEnv()

	viper.SetDefault("GS_DIR", "/opt/aimpanel/gs/")
	GS_DIR = viper.GetString("GS_DIR")

	viper.SetDefault("STORAGE_DIR", "/opt/aimpanel/storage/")
	STORAGE_DIR = viper.GetString("STORAGE_DIR")

	viper.SetDefault("TRASH_DIR", "/opt/aimpanel/trash/")
	TRASH_DIR = viper.GetString("TRASH_DIR")

	viper.SetDefault("BACKUP_DIR", "/opt/aimpanel/backups/")
	BACKUP_DIR = viper.GetString("BACKUP_DIR")

	viper.SetDefault("HOST_TOKEN", "")
	HOST_TOKEN = viper.GetString("HOST_TOKEN")

	viper.SetDefault("REDIS_HOST", "/opt/aimpanel/redis/redis.sock")
	REDIS_HOST = viper.GetString("REDIS_HOST")

	viper.SetDefault("REDIS_USERNAME", "")
	REDIS_USERNAME = viper.GetString("REDIS_USERNAME")

	viper.SetDefault("REDIS_PASSWORD", "")
	REDIS_PASSWORD = viper.GetString("REDIS_PASSWORD")

	viper.SetDefault("REDIS_PUB_SUB_AGENT_CH", "aimpanel_agent")
	REDIS_PUB_SUB_AGENT_CH = viper.GetString("REDIS_PUB_SUB_AGENT_CH")

	viper.SetDefault("REDIS_PUB_SUB_WRAPPER_CH", "aimpanel_wrapper")
	REDIS_PUB_SUB_WRAPPER_CH = viper.GetString("REDIS_PUB_SUB_WRAPPER_CH")

	viper.SetDefault("REDIS_PUB_SUB_SUPERVISOR_CH", "aimpanel_supervisor")
	REDIS_PUB_SUB_SUPERVISOR_CH = viper.GetString("REDIS_PUB_SUB_SUPERVISOR_CH")

	// use multiple URLs by separating them with space
	viper.SetDefault("MASTER_URLS", "https://aimpanel.local:3000")
	MASTER_URLS = viper.GetStringSlice("MASTER_URLS")

	viper.SetDefault("HW_ID", "")
	HW_ID = viper.GetString("HW_ID")
}
