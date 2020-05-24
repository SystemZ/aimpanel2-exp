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

	err = model.PutMetric(model.HostMetric, host.ID, metric.CpuUsage, taskMsg.CpuUsage, time.Now().Unix())
	if err != nil {
		return err
	}

	err = model.PutMetric(model.HostMetric, host.ID, metric.RamFree, taskMsg.RamFree, time.Now().Unix())
	if err != nil {
		return err
	}

	err = model.PutMetric(model.HostMetric, host.ID, metric.RamTotal, taskMsg.RamTotal, time.Now().Unix())
	if err != nil {
		return err
	}

	err = model.PutMetric(model.HostMetric, host.ID, metric.DiskFree, taskMsg.DiskFree, time.Now().Unix())
	if err != nil {
		return err
	}

	err = model.PutMetric(model.HostMetric, host.ID, metric.DiskUsed, taskMsg.DiskUsed, time.Now().Unix())
	if err != nil {
		return err
	}

	err = model.PutMetric(model.HostMetric, host.ID, metric.DiskTotal, taskMsg.DiskTotal, time.Now().Unix())
	if err != nil {
		return err
	}

	err = model.PutMetric(model.HostMetric, host.ID, metric.User, taskMsg.User, time.Now().Unix())
	if err != nil {
		return err
	}

	err = model.PutMetric(model.HostMetric, host.ID, metric.System, taskMsg.System, time.Now().Unix())
	if err != nil {
		return err
	}

	err = model.PutMetric(model.HostMetric, host.ID, metric.Idle, taskMsg.Idle, time.Now().Unix())
	if err != nil {
		return err
	}

	err = model.PutMetric(model.HostMetric, host.ID, metric.Nice, taskMsg.Nice, time.Now().Unix())
	if err != nil {
		return err
	}

	err = model.PutMetric(model.HostMetric, host.ID, metric.Iowait, taskMsg.Iowait, time.Now().Unix())
	if err != nil {
		return err
	}

	err = model.PutMetric(model.HostMetric, host.ID, metric.Irq, taskMsg.Irq, time.Now().Unix())
	if err != nil {
		return err
	}

	err = model.PutMetric(model.HostMetric, host.ID, metric.Softirq, taskMsg.Softirq, time.Now().Unix())
	if err != nil {
		return err
	}

	err = model.PutMetric(model.HostMetric, host.ID, metric.Steal, taskMsg.Steal, time.Now().Unix())
	if err != nil {
		return err
	}

	err = model.PutMetric(model.HostMetric, host.ID, metric.Guest, taskMsg.Guest, time.Now().Unix())
	if err != nil {
		return err
	}

	err = model.PutMetric(model.HostMetric, host.ID, metric.GuestNice, taskMsg.GuestNice, time.Now().Unix())
	if err != nil {
		return err
	}

	return nil
}
