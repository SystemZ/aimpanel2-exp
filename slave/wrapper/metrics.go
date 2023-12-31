package wrapper

import (
	proc "github.com/shirou/gopsutil/process"
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/lib/task"
	"gitlab.com/systemz/aimpanel2/slave/config"
	"gitlab.com/systemz/aimpanel2/slave/model"
	"time"
)

func (p *Process) Metrics() {
	for {
		<-time.After(time.Duration(p.MetricFrequency) * time.Second)

		if p.Running {
			process, err := proc.NewProcess(int32(p.Cmd.Process.Pid))
			if err != nil {
				logrus.Error(err.Error())
			}

			memoryInfoStat, err := process.MemoryInfo()
			if err != nil {
				logrus.Error(err.Error())
			}

			cpuPercent, err := process.CPUPercent()
			if err != nil {
				logrus.Error(err.Error())
			}

			rss := memoryInfoStat.RSS / 1024 / 1024

			taskMsg := task.Message{
				TaskId:       task.GAME_METRICS,
				GameServerID: p.GameServerID,
				CpuUsage:     int(cpuPercent),
				RamUsage:     int(rss),
			}

			model.SendTask(config.REDIS_PUB_SUB_AGENT_CH, taskMsg)
		}

	}
}
