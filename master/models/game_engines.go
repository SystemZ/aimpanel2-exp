package models

import (
	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

// Game engine represents the game engine for this application
// swagger:model game_engine
type GameEngine struct {
	// ID of the game
	//
	// required: true
	ID uint `gorm:"primary_key" json:"id"`

	// Game ID
	//
	// required: true
	GameId uuid.UUID `gorm:"column:game_id" json:"game_id"`

	// Order
	//
	// required: true
	Order int `gorm:"column:order" json:"order"`

	// Type
	//
	// required: true
	Type string `gorm:"column:type" json:"type"`

	// Version
	//
	// required: true
	Version string `gorm:"column:version" json:"version"`

	// Download URL
	//
	// required: true
	DownloadUrl string `gorm:"column:download_url" json:"download_url"`

	//Created at timestamp
	CreatedAt time.Time `json:"created_at"`

	//Updated at timestamp
	UpdatedAt time.Time `json:"updated_at"`

	//Deleted at timestamp
	DeletedAt *time.Time `json:"deleted_at"`
}

func (u *GameEngine) BeforeCreate(scope *gorm.Scope) error {
	uuidGen, err := uuid.NewV4()
	if err != nil {
		log.Println(err)
	}
	scope.SetColumn("ID", uuidGen)
	return nil
}
