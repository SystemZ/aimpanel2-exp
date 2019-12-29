package model

import (
	"github.com/gofrs/uuid"
	"time"
)

type Event struct {
	ID uint `gorm:"primary_key" json:"id"`

	EventId int `gorm:"column:event_id" json:"event_id"`

	HostId uuid.UUID `gorm:"column:host_id" json:"host_id"`

	UserId uuid.UUID `gorm:"column:user_id" json:"user_id"`

	GameServerId uuid.UUID `gorm:"column:game_server_id" json:"game_server_id"`

	//Created at timestamp
	CreatedAt time.Time `json:"created_at"`
}
