package gameserver

import (
	"encoding/json"
	"errors"
	"github.com/alexandrevicenzi/go-sse"
	"github.com/go-redis/redis"
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/lib/game"
	"gitlab.com/systemz/aimpanel2/lib/rabbit"
	"gitlab.com/systemz/aimpanel2/lib/task"
	"gitlab.com/systemz/aimpanel2/master/events"
	"gitlab.com/systemz/aimpanel2/master/model"
	"strconv"
	"time"
)

func HostData(hostToken string, taskMsg *task.Message) error {
	host := model.GetHostByToken(model.DB, hostToken)
	if host == nil {
		return errors.New("error when getting host from db")
	}

	switch taskMsg.TaskId {
	case rabbit.AGENT_METRICS_FREQUENCY:
		agentToken := taskMsg.AgentToken

		var host model.Host
		if model.DB.Where("token = ?", agentToken).First(&host).RecordNotFound() {
			break
		}

		channel, ok := events.SSE.GetChannel("/v1/events/" + hostToken)
		if !ok {
			return errors.New("game server is not turned on")
		}

		taskMsg := task.Message{
			TaskId:          task.AGENT_METRICS_FREQUENCY,
			MetricFrequency: host.MetricFrequency,
		}

		taskMsgStr, err := taskMsg.Serialize()
		if err != nil {
			return err
		}

		channel.SendMessage(sse.NewMessage("", taskMsgStr, strconv.Itoa(task.AGENT_METRICS_FREQUENCY)))
	case rabbit.AGENT_METRICS:
		var host model.Host
		if model.DB.Where("token = ?", taskMsg.AgentToken).First(&host).RecordNotFound() {
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
	case rabbit.AGENT_OS:
		var host model.Host
		if model.DB.Where("token = ?", taskMsg.AgentToken).First(&host).RecordNotFound() {
			break
		}

		host.OS = taskMsg.OS
		host.Platform = taskMsg.Platform
		host.PlatformFamily = taskMsg.PlatformFamily
		host.PlatformVersion = taskMsg.PlatformVersion
		host.KernelVersion = taskMsg.KernelVersion
		host.KernelArch = taskMsg.KernelArch

		model.DB.Save(&host)

	case rabbit.AGENT_HEARTBEAT:
		model.Redis.Set("agent_heartbeat_token_"+taskMsg.AgentToken, taskMsg.Timestamp, 24*time.Hour)
	}

	return nil
}

func GsData(hostToken string, gsId string, taskMsg *task.Message) error {
	switch taskMsg.TaskId {
	case task.WRAPPER_STARTED:
		logrus.Info("WRAPPER_STARTED")
		gameServerId := taskMsg.GameServerID
		_, err := model.Redis.Get("gs_restart_id_" + gameServerId).Int64()
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

			model.Redis.Del("gs_restart_id_" + gs.ID.String())
		}

		_, err = model.Redis.Get("gs_start_id_" + gameServerId).Int64()
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

			model.Redis.Del("gs_start_id_" + gs.ID.String())
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

		//sse.SSE.SendMessage("/v1/console/" + gsLog.GameServerId.String(),
		//	sse2.SimpleMessage(base64.StdEncoding.EncodeToString([]byte(gsLog.Log))))

		err := model.DB.Save(&gsLog).Error
		if err != nil {
			logrus.Warn(err)
		}
	case task.WRAPPER_EXITED:
		logrus.Info("WRAPPER_EXITED")
		gameServerId := taskMsg.GameServerID

		val, err := model.Redis.Get("gs_restart_id_" + gameServerId).Int()
		if err != redis.Nil && val != -1 {
			model.Redis.Set("gs_restart_id_"+gameServerId, 2, 24*time.Hour)

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

			model.Redis.Set("gs_restart_id_"+gameServerId, 3, 24*time.Hour)
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
		model.Redis.Set("wrapper_heartbeat_id_"+taskMsg.GameServerID, taskMsg.Timestamp, 24*time.Hour)
	}

	return nil
}
