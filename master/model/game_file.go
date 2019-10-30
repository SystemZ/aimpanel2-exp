package model

import (
	"github.com/jinzhu/gorm"
	"time"
)

type GameFile struct {
	ID uint `gorm:"primary_key" json:"id"`

	GameId uint `gorm:"column:game_id" json:"game_id"`

	GameVersion string `gorm:"column:game_version" json:"game_version"`

	DownloadUrl string `gorm:"column:download_url" json:"download_url"`

	//Created at timestamp
	CreatedAt time.Time `json:"created_at"`

	//Updated at timestamp
	UpdatedAt time.Time `json:"updated_at"`

	//Deleted at timestamp
	DeletedAt *time.Time `json:"deleted_at"`
}

func GetGameFileByGameIdAndVersion(db *gorm.DB, gameId uint, version string) *GameFile {
	var gf GameFile

	if db.Where("game_id = ? and (game_version = ? OR game_version = '0')", gameId, version).First(&gf).RecordNotFound() {
		return nil
	}

	return &gf
}
