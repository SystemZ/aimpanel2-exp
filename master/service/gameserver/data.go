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
)

func HostData(hostToken string, taskMsg *task.Message) error {
	logrus.Info("HostData")
	host := model.GetHostByToken(hostToken)
	if host == nil {
		return errors.New("error when getting host from db")
	}

	switch taskMsg.TaskId {
	case task.AGENT_STARTED:
		//model.DB.Save(&model.Event{
		//	EventId: event.AGENT_START,
		//	HostId:  host.ID,
		//})

		err := Update(host.ID)
		if err != nil {
			logrus.Error(err)
		}
	case task.AGENT_METRICS_FREQUENCY:
		logrus.Info("AGENT_METRICS_FREQUENCY")

		host := model.GetHostByToken(hostToken)
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
		host := model.GetHostByToken(hostToken)
		if host == nil {
			break
		}

		metric := &model.MetricHost{
			Base: model.Base{
				DocType: "metric_host",
			},
			HostId:    host.ID,
			CpuUsage:  taskMsg.CpuUsage,
			RamFree:   taskMsg.RamFree,
			RamTotal:  taskMsg.RamTotal,
			DiskFree:  taskMsg.DiskFree,
			DiskUsed:  taskMsg.DiskUsed,
			DiskTotal: taskMsg.DiskTotal,
			User:      taskMsg.User,
			System:    taskMsg.System,
			Idle:      taskMsg.Idle,
			Nice:      taskMsg.Nice,
			Iowait:    taskMsg.Iowait,
			Irq:       taskMsg.Irq,
			Softirq:   taskMsg.Softirq,
			Steal:     taskMsg.Steal,
			Guest:     taskMsg.Guest,
			GuestNice: taskMsg.GuestNice,
		}
		err := metric.Put(&metric)
		if err != nil {
			return err
		}
	case task.AGENT_OS:
		host := model.GetHostByToken(hostToken)
		if host == nil {
			break
		}

		host.OS = taskMsg.OS
		host.Platform = taskMsg.Platform
		host.PlatformFamily = taskMsg.PlatformFamily
		host.PlatformVersion = taskMsg.PlatformVersion
		host.KernelVersion = taskMsg.KernelVersion
		host.KernelArch = taskMsg.KernelArch

		err := host.Update(&host)
		if err != nil {
			return err
		}
	case task.AGENT_SHUTDOWN:
		//model.DB.Save(&model.Event{
		//	EventId: event.AGENT_SHUTDOWN,
		//	HostId:  host.ID,
		//})
	case task.AGENT_FILE_LIST_GS:
		logrus.Info("GAME_FILE_LIST")
		err := model.GsFilesPublish(model.Redis, taskMsg.GameServerID, &taskMsg.Files)
		if err != nil {
			logrus.Error(err)
		}
	}

	return nil
}

func GsData(hostToken string, taskMsg *task.Message) error {
	switch taskMsg.TaskId {
	case task.GAME_STARTED:
		logrus.Info("GAME_STARTED")

	case task.GAME_SERVER_LOG:
		logrus.Info("SERVER_LOG")
		var gsLog model.GameServerLog
		gsLog.Base.DocType = "game_server_log"
		gsLog.GameServerId = taskMsg.GameServerID

		if len(taskMsg.Stdout) > 0 {
			gsLog.Log = taskMsg.Stdout
			gsLog.Type = model.STDOUT
		}

		if len(taskMsg.Stderr) > 0 {
			gsLog.Log = taskMsg.Stderr
			gsLog.Type = model.STDERR
		}

		host := model.GetHostByToken(hostToken)
		events.SSE.SendMessage(fmt.Sprintf("/v1/host/%s/server/%s/console",
			host.ID,
			gsLog.GameServerId),
			sse.SimpleMessage(base64.StdEncoding.EncodeToString([]byte(gsLog.Log))))

		err := gsLog.Put(&gsLog)
		if err != nil {
			logrus.Warn(err)
		}
	case task.GAME_SHUTDOWN:
		logrus.Info("GAME_EXITED")

	case task.GAME_METRICS_FREQUENCY:
		gameServerId := taskMsg.GameServerID

		gs := model.GetGameServer(gameServerId)
		if gs == nil {
			break
		}

		channel, ok := events.SSE.GetChannel("/v1/events/" + hostToken)
		if !ok {
			return errors.New("game server is not turned on")
		}

		taskMsg := task.Message{
			TaskId:          task.GAME_METRICS_FREQUENCY,
			MetricFrequency: gs.MetricFrequency,
		}

		taskMsgStr, err := taskMsg.Serialize()
		if err != nil {
			return err
		}

		channel.SendMessage(sse.NewMessage("", taskMsgStr, taskMsg.TaskId.StringValue()))
	case task.GAME_METRICS:
		metric := &model.MetricGameServer{
			Base: model.Base{
				DocType: "metric_game_server",
			},
			GameServerId: taskMsg.GameServerID,
			CpuUsage:     taskMsg.CpuUsage,
			RamUsage:     taskMsg.RamUsage,
		}
		err := metric.Put(&metric)
		if err != nil {
			return err
		}
	}

	return nil
}
