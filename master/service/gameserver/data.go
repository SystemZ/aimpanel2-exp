package gameserver

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/alexandrevicenzi/go-sse"
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/lib/task"
	"gitlab.com/systemz/aimpanel2/master/events"
	"gitlab.com/systemz/aimpanel2/master/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func HostData(hostToken string, taskMsg *task.Message) error {
	host, err := model.GetHostByToken(hostToken)
	if err != nil {
		return err
	}

	if host == nil {
		return errors.New("error when getting host from db")
	}

	switch taskMsg.TaskId {
	case task.AGENT_STARTED:
		logrus.Infof("Got %v", taskMsg.TaskId)
		//model.DB.Save(&model.Event{
		//	EventId: event.AGENT_START,
		//	HostId:  host.ID,
		//})

		err := Update(host.ID)
		if err != nil {
			logrus.Error(err)
		}
	case task.AGENT_METRICS_FREQUENCY:
		logrus.Infof("Got %v", taskMsg.TaskId)

		host, err := model.GetHostByToken(hostToken)
		if err != nil {
			return err
		}

		if host == nil {
			break
		}

		channel, ok := events.SSE.GetChannel("/v1/events/" + hostToken)
		if !ok {
			return errors.New("host is not turned on")
		}

		taskMsg := task.Message{
			TaskId:          task.AGENT_METRICS_FREQUENCY,
			MetricFrequency: host.MetricFrequency,
		}

		taskMsgStr, err := taskMsg.Serialize()
		if err != nil {
			return err
		}

		channel.SendMessage(sse.NewMessage("", taskMsgStr, taskMsg.TaskId.StringValue()))
	case task.AGENT_METRICS:
		logrus.Infof("Got %v", taskMsg.TaskId)

		err := AgentMetrics(hostToken, *taskMsg)
		if err != nil {
			return err
		}
	case task.AGENT_OS:
		logrus.Infof("Got %v", taskMsg.TaskId)
		host, err := model.GetHostByToken(hostToken)
		if err != nil {
			return err
		}

		if host == nil {
			break
		}

		host.OS = taskMsg.OS
		host.Platform = taskMsg.Platform
		host.PlatformFamily = taskMsg.PlatformFamily
		host.PlatformVersion = taskMsg.PlatformVersion
		host.KernelVersion = taskMsg.KernelVersion
		host.KernelArch = taskMsg.KernelArch

		err = model.Update(host)
		if err != nil {
			return err
		}

	case task.AGENT_SHUTDOWN:
		logrus.Infof("Got %v", taskMsg.TaskId)
		//model.DB.Save(&model.Event{
		//	EventId: event.AGENT_SHUTDOWN,
		//	HostId:  host.ID,
		//})

	case task.AGENT_GET_JOBS:
		logrus.Infof("Got %v", taskMsg.TaskId)

		host, err := model.GetHostByToken(hostToken)
		if err != nil {
			return err
		}

		if host == nil {
			break
		}

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

		channel, ok := events.SSE.GetChannel("/v1/events/" + hostToken)
		if !ok {
			return errors.New("host is not turned on")
		}

		taskMsg := task.Message{
			TaskId: task.AGENT_GET_JOBS,
			Jobs:   &jobs,
		}

		taskMsgStr, err := taskMsg.Serialize()
		if err != nil {
			return err
		}

		channel.SendMessage(sse.NewMessage("", taskMsgStr, taskMsg.TaskId.StringValue()))

	default:
		logrus.Infof("Unhandled task %v", taskMsg.TaskId)
	}

	return nil
}

func GsData(hostToken string, taskMsg *task.Message) error {
	switch taskMsg.TaskId {
	case task.GAME_STARTED:
		logrus.Infof("Got %v", taskMsg.TaskId)

	case task.GAME_SERVER_LOG:
		logrus.Infof("Got %v", taskMsg.TaskId)
		var gsLog model.GameServerLog
		oid, _ := primitive.ObjectIDFromHex(taskMsg.GameServerID)
		gsLog.GameServerId = oid

		if len(taskMsg.Stdout) > 0 {
			gsLog.Log = taskMsg.Stdout
			gsLog.Type = model.STDOUT
		}

		if len(taskMsg.Stderr) > 0 {
			gsLog.Log = taskMsg.Stderr
			gsLog.Type = model.STDERR
		}

		host, err := model.GetHostByToken(hostToken)
		if err != nil {
			return err
		}

		events.SSE.SendMessage(fmt.Sprintf("/v1/host/%s/server/%s/console",
			host.ID.Hex(),
			gsLog.GameServerId.Hex()),
			sse.SimpleMessage(base64.StdEncoding.EncodeToString([]byte(gsLog.Log))))

		err = model.Put(&gsLog)
		if err != nil {
			logrus.Warn(err)
		}
	case task.GAME_SHUTDOWN:
		logrus.Infof("Got %v", taskMsg.TaskId)

	case task.GAME_METRICS_FREQUENCY:
		logrus.Infof("Got %v", taskMsg.TaskId)
		gameServerId := taskMsg.GameServerID

		oid, err := primitive.ObjectIDFromHex(taskMsg.GameServerID)
		if err != nil {
			return err
		}

		gs, err := model.GetGameServerById(oid)
		if err != nil {
			return err
		}

		if gs == nil {
			break
		}

		channel, ok := events.SSE.GetChannel("/v1/events/" + hostToken)
		if !ok {
			return errors.New("game server is not turned on")
		}

		taskMsg := task.Message{
			TaskId:          task.GAME_METRICS_FREQUENCY,
			GameServerID:    gameServerId,
			MetricFrequency: gs.MetricFrequency,
		}

		taskMsgStr, err := taskMsg.Serialize()
		if err != nil {
			return err
		}

		channel.SendMessage(sse.NewMessage("", taskMsgStr, taskMsg.TaskId.StringValue()))
	case task.GAME_METRICS:
		logrus.Infof("Got %v", taskMsg.TaskId)
		err := GameServerMetrics(*taskMsg)
		if err != nil {
			return err
		}

	case task.AGENT_FILE_LIST_GS:
		logrus.Infof("Got %v", taskMsg.TaskId)
		err := model.GsFilesPublish(model.Redis, taskMsg.GameServerID, taskMsg.Files)
		if err != nil {
			logrus.Error(err)
		}
	default:
		logrus.Infof("Unhandled task %v", taskMsg.TaskId)
	}

	return nil
}
