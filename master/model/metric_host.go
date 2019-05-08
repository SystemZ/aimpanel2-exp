package model

import (
	"github.com/gofrs/uuid"
	"time"
)

type MetricHost struct {
	ID uint `gorm:"primary_key" json:"id"`

	HostId uuid.UUID `gorm:"column:host_id" json:"host_id"`

	CpuUsage int `gorm:"column:cpu_usage" json:"cpu_usage"`

	RamFree int `gorm:"column:ram_free" json:"ram_free"`

	DiskFree int `gorm:"column:disk_free" json:"disk_free"`

	User      int `gorm:"column:user" json:"user"`

	System    int `gorm:"column:system" json:"system"`

	Idle      int `gorm:"column:idle" json:"idle"`

	Nice      int `gorm:"column:nice" json:"nice"`

	Iowait    int `gorm:"column:iowait" json:"iowait"`

	Irq       int `gorm:"column:irq" json:"irq"`

	Softirq   int `gorm:"column:softirq" json:"softirq"`

	Steal     int `gorm:"column:steal" json:"steal"`

	Guest     int `gorm:"column:guest" json:"guest"`

	GuestNice int `gorm:"column:guest_nice" json:"guest_nice"`

	Stolen    int `gorm:"column:stolen" json:"stolen"`

	//Created at timestamp
	CreatedAt time.Time `json:"created_at"`
}

func (m *MetricHost) TableName() string {
	return "metrics_host"
}
