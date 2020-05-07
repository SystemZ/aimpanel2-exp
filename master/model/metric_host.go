package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type MetricHost struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty" example:"1238206236281802752"`

	// ID of the host
	HostId primitive.ObjectID `bson:"host_id" json:"host_id" example:"100112233-4455-6677-8899-aabbccddeeff"`

	// CPU usage
	CpuUsage int `bson:"cpu_usage" json:"cpu_usage" example:"12"`

	// Ram free
	RamFree int `bson:"ram_free" json:"ram_free" example:"4654"`

	// Ram total
	RamTotal int `bson:"ram_total" json:"ram_total" example:"15591"`

	// Disk free
	DiskFree int `bson:"disk_free" json:"disk_free" example:"221024"`

	// Disk used
	DiskUsed int `bson:"disk_used" json:"disk_used" example:"19748"`

	// Disk total
	DiskTotal int `bson:"disk_total" json:"disk_total" example:"253730"`

	// User
	User int `bson:"user" json:"user" example:"6400"`

	// System
	System int `bson:"system" json:"system" example:"1546"`

	// Idle
	Idle int `bson:"idle" json:"idle" example:"52860"`

	// Nice
	Nice int `bson:"nice" json:"nice" example:"6"`

	// IO Wait
	Iowait int `bson:"iowait" json:"iowait" example:"32"`

	// Irq
	Irq int `bson:"irq" json:"irq" example:"0"`

	// Soft irq
	Softirq int `bson:"softirq" json:"softirq" example:"703"`

	// Steal
	Steal int `bson:"steal" json:"steal" example:"0"`

	// Guest
	Guest int `bson:"guest" json:"guest" example:"0"`

	// Guest nice
	GuestNice int `bson:"guest_nice" json:"guest_nice" example:"0"`
}

func (m *MetricHost) GetCollectionName() string {
	return metricHostCollection
}

func (m *MetricHost) GetID() primitive.ObjectID {
	return m.ID
}

func (m *MetricHost) SetID(id primitive.ObjectID) {
	m.ID = id
}
