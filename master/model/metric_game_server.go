package model

import (
	"github.com/gofrs/uuid"
	"time"
)

type MetricGameServer struct {
	ID uint `gorm:"primary_key" json:"id"`

	GameServerId uuid.UUID `gorm:"column:game_server_id" json:"game_server_id"`

	RamUsage int `gorm:"column:ram_usage" json:"ram_usage"`

	CpuUsage int `gorm:"column:cpu_usage" json:"cpu_usage"`

	//Created at timestamp
	CreatedAt time.Time `json:"created_at"`
}

func (m *MetricGameServer) TableName() string {
	return "metrics_game_server"
}