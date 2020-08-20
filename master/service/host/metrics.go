package host

import (
	"gitlab.com/systemz/aimpanel2/lib/metric"
	"gitlab.com/systemz/aimpanel2/lib/task"
	"gitlab.com/systemz/aimpanel2/master/model"
	"time"
)

func Metrics(hostToken string, taskMsg task.Message) error {
	host, err := model.GetHostByToken(hostToken)
	if err != nil {
		return err
	}

	if host == nil {
		return nil
	}

	type mList struct {
		metricId       metric.Id
		taskPropertyId int
	}
	metricList := []mList{
		{metricId: metric.RamFree, taskPropertyId: taskMsg.RamFree},
		{metricId: metric.CpuUsage, taskPropertyId: taskMsg.CpuUsage},
		{metricId: metric.RamTotal, taskPropertyId: taskMsg.RamTotal},
		{metricId: metric.DiskFree, taskPropertyId: taskMsg.DiskFree},
		{metricId: metric.DiskUsed, taskPropertyId: taskMsg.DiskUsed},
		{metricId: metric.DiskTotal, taskPropertyId: taskMsg.DiskTotal},
		{metricId: metric.User, taskPropertyId: taskMsg.User},
		{metricId: metric.System, taskPropertyId: taskMsg.System},
		{metricId: metric.Idle, taskPropertyId: taskMsg.Idle},
		{metricId: metric.Nice, taskPropertyId: taskMsg.Nice},
		{metricId: metric.Iowait, taskPropertyId: taskMsg.Iowait},
		{metricId: metric.Irq, taskPropertyId: taskMsg.Irq},
		{metricId: metric.Softirq, taskPropertyId: taskMsg.Softirq},
		{metricId: metric.Steal, taskPropertyId: taskMsg.Steal},
		{metricId: metric.Guest, taskPropertyId: taskMsg.Guest},
		{metricId: metric.GuestNice, taskPropertyId: taskMsg.GuestNice},
	}

	// loop over all metrics for write in DB
	for _, v := range metricList {
		err = model.PutMetric(model.HostMetric, host.ID, v.metricId, v.taskPropertyId, time.Now().Unix())
		if err != nil {
			return err
		}
	}

	return nil
}
