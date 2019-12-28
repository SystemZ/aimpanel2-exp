package model

import (
	"github.com/gofrs/uuid"
	"time"
)

type MetricHost struct {
	// ID of the metric
	ID uint `gorm:"primary_key" json:"id" example:"1337"`

	// ID of the host
	HostId uuid.UUID `gorm:"column:host_id" json:"host_id" example:"100112233-4455-6677-8899-aabbccddeeff"`

	// CPU usage
	CpuUsage int `gorm:"column:cpu_usage" json:"cpu_usage" example:"12"`

	// Ram free
	RamFree int `gorm:"column:ram_free" json:"ram_free" example:"4654"`

	// Ram total
	RamTotal int `gorm:"column:ram_total" json:"ram_total" example:"15591"`

	// Disk free
	DiskFree int `gorm:"column:disk_free" json:"disk_free" example:"221024"`

	// Disk used
	DiskUsed int `gorm:"column:disk_used" json:"disk_used" example:"19748"`

	// Disk total
	DiskTotal int `gorm:"column:disk_total" json:"disk_total" example:"253730"`

	// User
	User int `gorm:"column:user" json:"user" example:"6400"`

	// System
	System int `gorm:"column:system" json:"system" example:"1546"`

	// Idle
	Idle int `gorm:"column:idle" json:"idle" example:"52860"`

	// Nice
	Nice int `gorm:"column:nice" json:"nice" example:"6"`

	// IO Wait
	Iowait int `gorm:"column:iowait" json:"iowait" example:"32"`

	// Irq
	Irq int `gorm:"column:irq" json:"irq" example:"0"`

	// Soft irq
	Softirq int `gorm:"column:softirq" json:"softirq" example:"703"`

	// Steal
	Steal int `gorm:"column:steal" json:"steal" example:"0"`

	// Guest
	Guest int `gorm:"column:guest" json:"guest" example:"0"`

	// Guest nice
	GuestNice int `gorm:"column:guest_nice" json:"guest_nice" example:"0"`

	// Date of metric creation
	CreatedAt time.Time `json:"created_at" example:"2019-09-29T03:16:27+02:00"`
}

func (m *MetricHost) TableName() string {
	return "metrics_host"
}
