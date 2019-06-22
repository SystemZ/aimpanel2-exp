package model

import (
	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	"time"
)

const (
	STDOUT = iota
	STDERR = iota
)

type GameServerLog struct {
	ID uint `gorm:"primary_key" json:"id"`

	GameServerId uuid.UUID `gorm:"column:game_server_id" json:"game_server_id"`

	Type uint `gorm:"column:type" json:"type"`

	Log string `gorm:"column:log" json:"log"`

	//Created at timestamp
	CreatedAt time.Time `json:"created_at"`
}

func GetLogsByGameServer(db *gorm.DB, gsId string, limit int) *[]GameServerLog {
	var logs []GameServerLog

	if db.Where("game_server_id = ?", gsId).Order("created_at desc").Limit(limit).Find(&logs).RecordNotFound() {
		return nil
	}

	return &logs
}
