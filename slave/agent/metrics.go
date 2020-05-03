package agent

import (
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/lib/ahttp"
	"gitlab.com/systemz/aimpanel2/lib/task"
	"gitlab.com/systemz/aimpanel2/slave/config"
	"time"
)

func metrics() {
	for {
		<-time.After(time.Duration(metricsFrequency) * time.Second)

		virtualMemory, _ := mem.VirtualMemory()
		ramFree := virtualMemory.Free / 1024 / 1024
		ramTotal := virtualMemory.Total / 1024 / 1024
		cpuPercent, _ := cpu.Percent(time.Duration(1)*time.Second, false)
		cpuTimes, _ := cpu.Times(false)

		diskUsage, _ := disk.Usage("/")
		diskFree := diskUsage.Free / 1024 / 1024
		diskTotal := diskUsage.Total / 1024 / 1024
		diskUsed := diskUsage.Used / 1024 / 1024

		taskMsg := task.Message{
			TaskId:    task.AGENT_METRICS,
			CpuUsage:  int(cpuPercent[0]),
			RamFree:   int(ramFree),
			RamTotal:  int(ramTotal),
			DiskFree:  int(diskFree),
			DiskTotal: int(diskTotal),
			DiskUsed:  int(diskUsed),
			User:      int(cpuTimes[0].User),
			System:    int(cpuTimes[0].System),
			Idle:      int(cpuTimes[0].Idle),
			Nice:      int(cpuTimes[0].Nice),
			Iowait:    int(cpuTimes[0].Iowait),
			Irq:       int(cpuTimes[0].Irq),
			Softirq:   int(cpuTimes[0].Softirq),
			Steal:     int(cpuTimes[0].Steal),
			Guest:     int(cpuTimes[0].Guest),
			GuestNice: int(cpuTimes[0].GuestNice),
		}
		//TODO: do something with status code
		_, err := ahttp.SendTaskData(config.API_URL+"/v1/events/"+config.HOST_TOKEN, config.API_TOKEN, taskMsg)
		if err != nil {
			logrus.Error(err)
		}
	}
}

func sendOSInfo() {
	h, _ := host.Info()

	taskMsg := task.Message{
		TaskId: task.AGENT_OS,

		OS:              h.OS,
		Platform:        h.Platform,
		PlatformFamily:  h.PlatformFamily,
		PlatformVersion: h.PlatformVersion,
		KernelVersion:   h.KernelVersion,
		KernelArch:      h.KernelArch,
	}
	//TODO: do something with status code
	_, err := ahttp.SendTaskData(config.API_URL+"/v1/events/"+config.HOST_TOKEN, config.API_TOKEN, taskMsg)
	if err != nil {
		logrus.Error(err)
	}
}
