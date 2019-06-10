package model

import (
	"github.com/jinzhu/gorm"
	"time"
)

// Game represents the game for this application
// swagger:model game
type Game struct {
	// ID of the game
	//
	// required: true
	ID uint `gorm:"primary_key" json:"id"`

	// Name
	//
	// required: true
	Name string `gorm:"column:name" json:"name"`

	//Created at timestamp
	CreatedAt time.Time `json:"created_at"`

	//Updated at timestamp
	UpdatedAt time.Time `json:"updated_at"`

	//Deleted at timestamp
	DeletedAt *time.Time `json:"deleted_at"`
}

func GetGameStartCommandByVersion(db *gorm.DB, gameId uint, versionId uint) *GameCommand {
	var startCommand GameCommand

	if db.Where("game_id = ? and type = ? and game_version = ?", gameId, "start", versionId).
		First(&startCommand).RecordNotFound() {
		return nil
	}

	return &startCommand
}

func GetGameInstallCommandsByVersion(db *gorm.DB, gameId uint, versionId uint) *[]GameCommand {
	var installCommands []GameCommand

	if db.Where("game_id = ? and type = ? and game_version = ?", gameId, "install", versionId).
		Order("`order` asc").Find(&installCommands).RecordNotFound() {
		return nil
	}

	return &installCommands
}

func GetGameInstallFileByVersion(db *gorm.DB, gameId uint, versionId uint) *GameFile {
	var file GameFile

	if db.Where("game_id = ? and game_version = ?", gameId, versionId).First(&file).RecordNotFound() {
		return nil
	}

	return &file
}
