package model

type MetricHost struct {
	Base

	// ID of the host
	HostId string `json:"host_id" example:"100112233-4455-6677-8899-aabbccddeeff"`

	// CPU usage
	CpuUsage int `json:"cpu_usage" example:"12"`

	// Ram free
	RamFree int `json:"ram_free" example:"4654"`

	// Ram total
	RamTotal int `json:"ram_total" example:"15591"`

	// Disk free
	DiskFree int `json:"disk_free" example:"221024"`

	// Disk used
	DiskUsed int `json:"disk_used" example:"19748"`

	// Disk total
	DiskTotal int `json:"disk_total" example:"253730"`

	// User
	User int `json:"user" example:"6400"`

	// System
	System int `json:"system" example:"1546"`

	// Idle
	Idle int `json:"idle" example:"52860"`

	// Nice
	Nice int `json:"nice" example:"6"`

	// IO Wait
	Iowait int `json:"iowait" example:"32"`

	// Irq
	Irq int `json:"irq" example:"0"`

	// Soft irq
	Softirq int `json:"softirq" example:"703"`

	// Steal
	Steal int `json:"steal" example:"0"`

	// Guest
	Guest int `json:"guest" example:"0"`

	// Guest nice
	GuestNice int `json:"guest_nice" example:"0"`
}
