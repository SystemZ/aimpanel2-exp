package model

import "time"

type GameFile struct {
	ID uint `gorm:"primary_key" json:"id"`

	Type string `gorm:"column:type" json:"name"`

	GameId uint `gorm:"column:game_id" json:"game_id"`

	GameVersion uint `gorm:"column:game_version" json:"game_version"`

	DownloadUrl string `gorm:"column:download_url" json:"download_url"`

	Filename string `gorm:"column:filename" json:"filename"`

	//Created at timestamp
	CreatedAt time.Time `json:"created_at"`

	//Updated at timestamp
	UpdatedAt time.Time `json:"updated_at"`

	//Deleted at timestamp
	DeletedAt *time.Time `json:"deleted_at"`
}
