package model

import (
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
