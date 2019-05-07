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

func (g *Game) GetStartCommand(db *gorm.DB) *GameCommand {
	var startCommand GameCommand

	if db.Where("game_id = ? and type = ?", g.ID, "start").First(&startCommand).RecordNotFound() {
		return nil
	}

	return &startCommand
}

func (g *Game) GetStartCommandByVersion(db *gorm.DB, versionId uint) *GameCommand {
	var startCommand GameCommand

	if db.Where("game_id = ? and type = ? and game_version = ?", g.ID, "start", versionId).
		First(&startCommand).RecordNotFound() {
		return nil
	}

	return &startCommand
}

func (g *Game) GetInstallCommands(db *gorm.DB) *[]GameCommand {
	var installCommands []GameCommand

	if db.Where("game_id = ? and type = ?", g.ID, "install").
		Order("`order` asc").Find(&installCommands).RecordNotFound() {
		return nil
	}

	return &installCommands
}

func (g *Game) GetInstallCommandsByVersion(db *gorm.DB, versionId uint) *[]GameCommand {
	var installCommands []GameCommand

	if db.Where("game_id = ? and type = ? and game_version = ?", g.ID, "install", versionId).
		Order("`order` asc").Find(&installCommands).RecordNotFound() {
		return nil
	}

	return &installCommands
}

func (g *Game) GetInstallFile(db *gorm.DB) *GameFile {
	var file GameFile

	if db.Where("game_id = ?", g.ID).First(&file).RecordNotFound() {
		return nil
	}

	return &file
}

func (g *Game) GetInstallFileByVersion(db *gorm.DB, versionId uint) *GameFile {
	var file GameFile

	if db.Where("game_id = ? and game_version = ?", g.ID, versionId).First(&file).RecordNotFound() {
		return nil
	}

	return &file
}
