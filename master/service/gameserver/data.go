package gameserver

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alexandrevicenzi/go-sse"
	"github.com/go-redis/redis"
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/lib/game"
	"gitlab.com/systemz/aimpanel2/lib/task"
	"gitlab.com/systemz/aimpanel2/master/events"
	"gitlab.com/systemz/aimpanel2/master/model"
	"strconv"
)

func HostData(hostToken string, taskMsg *task.Message) error {
	logrus.Info("HostData")
	host := model.GetHostByToken(model.DB, hostToken)
	if host == nil {
		return errors.New("error when getting host from db")
	}

	switch taskMsg.TaskId {
	case task.AGENT_STARTED:
		//model.DB.Save(&model.Event{
		//	EventId: event.AGENT_START,
		//	HostId:  host.ID,
		//})

		err := Update(host.ID.String())
		if err != nil {
			logrus.Error(err)
		}
	case task.AGENT_METRICS_FREQUENCY:
		logrus.Info("AGENT_METRICS_FREQUENCY")

		var host model.Host
		if model.DB.Where("token = ?", hostToken).First(&host).RecordNotFound() {
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

		channel.SendMessage(sse.NewMessage("", taskMsgStr, strconv.Itoa(taskMsg.TaskId)))
	case task.AGENT_METRICS:
		var host model.Host
		if model.DB.Where("token = ?", hostToken).First(&host).RecordNotFound() {
			break
		}

		metric := &model.MetricHost{
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
		model.DB.Save(metric)
	case task.AGENT_OS:
		var host model.Host
		if model.DB.Where("token = ?", hostToken).First(&host).RecordNotFound() {
			break
		}

		host.OS = taskMsg.OS
		host.Platform = taskMsg.Platform
		host.PlatformFamily = taskMsg.PlatformFamily
		host.PlatformVersion = taskMsg.PlatformVersion
		host.KernelVersion = taskMsg.KernelVersion
		host.KernelArch = taskMsg.KernelArch

		model.DB.Save(&host)

	case task.AGENT_HEARTBEAT:
		model.SetAgentHeartbeat(model.Redis, hostToken, taskMsg.Timestamp)
	case task.AGENT_SHUTDOWN:
		//model.DB.Save(&model.Event{
		//	EventId: event.AGENT_SHUTDOWN,
		//	HostId:  host.ID,
		//})
	}

	return nil
}

func GsData(hostToken string, gsId string, taskMsg *task.Message) error {
	switch taskMsg.TaskId {
	case task.WRAPPER_STARTED:
		logrus.Info("WRAPPER_STARTED")
		gameServerId := taskMsg.GameServerID
		_, err := model.GetGsRestart(model.Redis, gameServerId)
		if err == nil {
			var gs model.GameServer
			if model.DB.Where("id = ?", gameServerId).First(&gs).RecordNotFound() {
				break
			}

			var gameDef game.Game
			err = json.Unmarshal([]byte(gs.GameJson), &gameDef)

			channel, ok := events.SSE.GetChannel("/v1/events/" + hostToken + "/" + gsId)
			if !ok {
				return errors.New("game server is not turned on")
			}

			taskMsg := task.Message{
				TaskId:       task.GAME_START,
				GameServerID: taskMsg.GameServerID,
				Game:         gameDef,
			}

			taskMsgStr, err := taskMsg.Serialize()
			if err != nil {
				return err
			}

			channel.SendMessage(sse.NewMessage("", taskMsgStr, strconv.Itoa(task.GAME_START)))

			model.DelGsRestart(model.Redis, gs.ID.String())
		}

		_, err = model.GetGsStart(model.Redis, gameServerId)
		if err == nil {
			var gs model.GameServer
			if model.DB.Where("id = ?", gameServerId).First(&gs).RecordNotFound() {
				break
			}

			var gameDef game.Game
			err = json.Unmarshal([]byte(gs.GameJson), &gameDef)

			channel, ok := events.SSE.GetChannel("/v1/events/" + hostToken + "/" + gsId)
			if !ok {
				return errors.New("game server is not turned on")
			}

			taskMsg := task.Message{
				TaskId:       task.GAME_START,
				GameServerID: taskMsg.GameServerID,
				Game:         gameDef,
			}

			taskMsgStr, err := taskMsg.Serialize()
			if err != nil {
				return err
			}

			channel.SendMessage(sse.NewMessage("", taskMsgStr, strconv.Itoa(task.GAME_START)))

			model.DelGsStart(model.Redis, gs.ID.String())
		}

	case task.SERVER_LOG:
		logrus.Info("SERVER_LOG")
		var gsLog model.GameServerLog
		gsLog.GameServerId = uuid.FromStringOrNil(taskMsg.GameServerID)

		if len(taskMsg.Stdout) > 0 {
			gsLog.Log = taskMsg.Stdout
			gsLog.Type = model.STDOUT
		}

		if len(taskMsg.Stderr) > 0 {
			gsLog.Log = taskMsg.Stderr
			gsLog.Type = model.STDERR
		}

		host := model.GetHostByToken(model.DB, hostToken)
		events.SSE.SendMessage(fmt.Sprintf("/v1/host/%s/server/%s/console",
			host.ID.String(),
			gsLog.GameServerId.String()),
			sse.SimpleMessage(base64.StdEncoding.EncodeToString([]byte(gsLog.Log))))

		err := model.DB.Save(&gsLog).Error
		if err != nil {
			logrus.Warn(err)
		}
	case task.WRAPPER_EXITED:
		logrus.Info("WRAPPER_EXITED")
		gameServerId := taskMsg.GameServerID

		val, err := model.GetGsRestart(model.Redis, gameServerId)
		if err != redis.Nil && val != -1 {
			model.SetGsRestart(model.Redis, gameServerId, 2)

			var gs model.GameServer
			if model.DB.Where("id = ?", gameServerId).First(&gs).RecordNotFound() {
				break
			}

			var host model.Host
			if model.DB.Where("id = ?", gs.HostId).First(&host).RecordNotFound() {
				break
			}

			channel, ok := events.SSE.GetChannel("/v1/events/" + hostToken + "/" + gsId)
			if !ok {
				return errors.New("game server is not turned on")
			}

			taskMsg := task.Message{
				TaskId:       task.WRAPPER_START,
				GameServerID: gameServerId,
			}

			taskMsgStr, err := taskMsg.Serialize()
			if err != nil {
				return err
			}

			channel.SendMessage(sse.NewMessage("", taskMsgStr, strconv.Itoa(task.WRAPPER_START)))

			model.SetGsRestart(model.Redis, gameServerId, 3)
		}

	case task.WRAPPER_METRICS_FREQUENCY:
		gameServerId := taskMsg.GameServerID

		var gs model.GameServer
		if model.DB.Where("id = ?", gameServerId).First(&gs).RecordNotFound() {
			break
		}

		channel, ok := events.SSE.GetChannel("/v1/events/" + hostToken + "/" + gsId)
		if !ok {
			return errors.New("game server is not turned on")
		}

		taskMsg := task.Message{
			TaskId:          task.WRAPPER_METRICS_FREQUENCY,
			MetricFrequency: gs.MetricFrequency,
		}

		taskMsgStr, err := taskMsg.Serialize()
		if err != nil {
			return err
		}

		channel.SendMessage(sse.NewMessage("", taskMsgStr, strconv.Itoa(task.WRAPPER_METRICS_FREQUENCY)))
	case task.WRAPPER_METRICS:
		metric := &model.MetricGameServer{
			GameServerId: uuid.FromStringOrNil(taskMsg.GameServerID),
			CpuUsage:     taskMsg.CpuUsage,
			RamUsage:     taskMsg.RamUsage,
		}
		model.DB.Save(metric)
	case task.WRAPPER_HEARTBEAT:
		model.SetWrapperHeartbeat(model.Redis, taskMsg.GameServerID, taskMsg.Timestamp)
	}

	return nil
}
