package host

import (
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/lib"
	"gitlab.com/systemz/aimpanel2/lib/ecode"
	"gitlab.com/systemz/aimpanel2/lib/task"
	"gitlab.com/systemz/aimpanel2/master/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Data(host model.Host, taskMsg *task.Message) error {
	switch taskMsg.TaskId {
	case task.AGENT_STARTED:
		logrus.Infof("Got %v", taskMsg.TaskId)

		// FIXME check for updates
		//err := AgentStarted(host.ID)
		//if err != nil {
		//	return err
		//}
	case task.AGENT_METRICS_FREQUENCY:
		logrus.Infof("Got %v", taskMsg.TaskId)

		err := AgentMetricsFrequency(&host)
		if err != nil {
			return err
		}
	case task.AGENT_METRICS:
		logrus.Infof("Got %v", taskMsg.TaskId)

		err := Metrics(host.Token, *taskMsg)
		if err != nil {
			return err
		}
	case task.AGENT_OS:
		logrus.Infof("Got %v", taskMsg.TaskId)
		err := AgentOS(&host, taskMsg)
		if err != nil {
			return err
		}

	case task.AGENT_SHUTDOWN:
		logrus.Infof("Got %v", taskMsg.TaskId)

	case task.AGENT_GET_JOBS:
		logrus.Infof("Got %v", taskMsg.TaskId)

		err := AgentGetJobs(&host)
		if err != nil {
			return err
		}
	default:
		logrus.Infof("Unhandled task %v", taskMsg.TaskId)
	}

	return nil
}

func AgentStarted(hostId primitive.ObjectID) error {
	err := Update(hostId, model.User{})
	if err != nil {
		return err
	}

	return nil
}

func AgentMetricsFrequency(host *model.Host) error {
	taskMsg := task.Message{
		TaskId:          task.AGENT_METRICS_FREQUENCY,
		MetricFrequency: host.MetricFrequency,
	}

	err := model.SendTaskToSlave(host.ID, model.User{}, taskMsg)
	if err != nil {
		return &lib.Error{ErrorCode: ecode.DbSave}
	}

	return nil
}

func AgentOS(host *model.Host, taskMsg *task.Message) error {
	host.OS = taskMsg.OS
	host.Platform = taskMsg.Platform
	host.PlatformFamily = taskMsg.PlatformFamily
	host.PlatformVersion = taskMsg.PlatformVersion
	host.KernelVersion = taskMsg.KernelVersion
	host.KernelArch = taskMsg.KernelArch

	err := model.Update(host)
	if err != nil {
		return err
	}

	return nil
}

func AgentGetJobs(host *model.Host) error {
	var jobs []task.Job

	hostJobs, err := model.GetHostJobsByHostId(host.ID)
	if err != nil {
		return err
	}

	for _, job := range hostJobs {
		jobs = append(jobs, task.Job{
			Name:           job.Name,
			CronExpression: job.CronExpression,
			TaskMessage:    job.TaskMessage,
		})
	}

	taskMsg := task.Message{
		TaskId: task.AGENT_GET_JOBS,
		Jobs:   &jobs,
	}

	err = model.SendTaskToSlave(host.ID, model.User{}, taskMsg)
	if err != nil {
		return &lib.Error{ErrorCode: ecode.DbSave}
	}

	return nil
}
