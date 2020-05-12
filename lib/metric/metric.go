package metric

import "strconv"

type Id int

func (i Id) StringValue() string {
	return strconv.Itoa(int(i))
}

const (
	CpuUsage Id = iota + 1
	RamUsage
	RamFree
	RamTotal
	DiskFree
	DiskUsed
	DiskTotal
	User
	System
	Idle
	Nice
	Iowait
	Irq
	Softirq
	Steal
	Guest
	GuestNice
)
