package model

import "time"

type GameCommand struct {
	ID uint `gorm:"primary_key" json:"id"`

	Type string `gorm:"column:type" json:"name"`

	GameId uint `gorm:"column:game_id" json:"game_id"`

	Version string `gorm:"column:version" json:"version"`

	Order uint `gorm:"column:order" json:"order"`

	Command string `gorm:"column:command" json:"command"`

	//Created at timestamp
	CreatedAt time.Time `json:"created_at"`

	//Updated at timestamp
	UpdatedAt time.Time `json:"updated_at"`

	//Deleted at timestamp
	DeletedAt *time.Time `json:"deleted_at"`
}
