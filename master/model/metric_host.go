package model

import (
	"github.com/gofrs/uuid"
	"time"
)

// Single metric
// swagger:model MetricHost
type MetricHost struct {
	// ID of the metric
	// Example: 1337
	ID uint `gorm:"primary_key" json:"id"`

	// ID of the host
	// Example: 100112233-4455-6677-8899-aabbccddeeff
	HostId uuid.UUID `gorm:"column:host_id" json:"host_id"`

	// CPU usage
	// Example: 12
	CpuUsage int `gorm:"column:cpu_usage" json:"cpu_usage"`

	// Ram free
	// Example: 4654
	RamFree int `gorm:"column:ram_free" json:"ram_free"`

	// Ram total
	// Example: 15591
	RamTotal int `gorm:"column:ram_total" json:"ram_total"`

	// Disk free
	// Example: 221024
	DiskFree int `gorm:"column:disk_free" json:"disk_free"`

	// Disk used
	// Example: 19748
	DiskUsed int `gorm:"column:disk_used" json:"disk_used"`

	// Disk total
	// Example: 253730
	DiskTotal int `gorm:"column:disk_total" json:"disk_total"`

	// User
	// Example: 6400
	User int `gorm:"column:user" json:"user"`

	// System
	// Example: 1546
	System int `gorm:"column:system" json:"system"`

	// Idle
	// Example: 52860
	Idle int `gorm:"column:idle" json:"idle"`

	// Nice
	// Example: 6
	Nice int `gorm:"column:nice" json:"nice"`

	// IO Wait
	// Example: 32
	Iowait int `gorm:"column:iowait" json:"iowait"`

	// Irq
	// Example: 0
	Irq int `gorm:"column:irq" json:"irq"`

	// Soft irq
	// Example: 703
	Softirq int `gorm:"column:softirq" json:"softirq"`

	// Steal
	// Example: 0
	Steal int `gorm:"column:steal" json:"steal"`

	// Guest
	// Example: 0
	Guest int `gorm:"column:guest" json:"guest"`

	// Guest nice
	// Example: 0
	GuestNice int `gorm:"column:guest_nice" json:"guest_nice"`

	// Date of metric creation
	// Example: 2019-09-29T03:16:27+02:00
	CreatedAt time.Time `json:"created_at"`
}

func (m *MetricHost) TableName() string {
	return "metrics_host"
}
