package model

import (
	"github.com/gofrs/uuid"
	"time"
)

type MetricHost struct {
	ID uint `gorm:"primary_key" json:"id"`

	HostId uuid.UUID `gorm:"column:host_id" json:"host_id"`

	CpuUsage int64 `gorm:"column:cpu_usage" json:"cpu_usage"`

	RamFree int64 `gorm:"column:ram_free" json:"ram_free"`

	DiskFree int64 `gorm:"column:disk_free" json:"disk_free"`

	//Created at timestamp
	CreatedAt time.Time `json:"created_at"`
}

func (m *MetricHost) TableName() string {
	return "metrics_host"
}