package model

import (
	"github.com/gofrs/uuid"
	"time"
)

const (
	STDOUT = iota
	STDERR = iota
)

type GameServerLog struct {
	ID uint `gorm:"primary_key" json:"id"`

	GameServerID uuid.UUID `gorm:"column:game_server_id" json:"game_server_id"`

	Type uint `gorm:"column:type" json:"type"`

	Log string `gorm:"column:log" json:"log"`

	//Created at timestamp
	CreatedAt time.Time `json:"created_at"`
}
