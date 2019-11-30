package rabbit

import (
	"encoding/base64"
	"encoding/json"
	sse2 "github.com/alexandrevicenzi/go-sse"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/lib"
	"gitlab.com/systemz/aimpanel2/lib/game"
	"gitlab.com/systemz/aimpanel2/lib/rabbit"
	"gitlab.com/systemz/aimpanel2/master/model"
	"gitlab.com/systemz/aimpanel2/master/sse"
	"time"
)

func ListenWrapperData() {
	queue, err := channel.QueueDeclare(
		"wrapper_data",
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	lib.FailOnError(err, "Failed to declare a queue")

	msgs, err := channel.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil)
	lib.FailOnError(err, "Failed to register a consumer")

	go func() {
		for msg := range msgs {
			var msgBody rabbit.QueueMsg
			err = json.Unmarshal(msg.Body, &msgBody)
			if err != nil {
				logrus.Warn(err)
			}

			switch msgBody.TaskId {
			case rabbit.SERVER_LOG:
				logrus.Info("SERVER_LOG")
				var gsLog model.GameServerLog
				gsLog.GameServerId = msgBody.GameServerID

				if len(msgBody.Stdout) > 0 {
					gsLog.Log = msgBody.Stdout
					gsLog.Type = model.STDOUT
				}

				if len(msgBody.Stderr) > 0 {
					gsLog.Log = msgBody.Stderr
					gsLog.Type = model.STDERR
				}

				sse.SSE.SendMessage("/v1/console/" + gsLog.GameServerId.String(),
					sse2.SimpleMessage(base64.StdEncoding.EncodeToString([]byte(gsLog.Log))))

				err = model.DB.Save(&gsLog).Error
				if err != nil {
					logrus.Warn(err)
				}
			case rabbit.WRAPPER_STARTED:
				logrus.Info("WRAPPER_STARTED")
				gameServerId := msgBody.GameServerID
				_, err := model.Redis.Get("gs_restart_id_" + gameServerId.String()).Int64()
				if err == nil {
					var gs model.GameServer
					if model.DB.Where("id = ?", gameServerId).First(&gs).RecordNotFound() {
						break
					}

					//var startCommand model.GameCommand
					//if model.DB.Where("game_id = ? and type = ?", gs.GameId, "start").
					//	First(&startCommand).RecordNotFound() {
					//	break
					//}
					var gameDef game.Game
					err = json.Unmarshal([]byte(gs.GameJson), &gameDef)

					msg := rabbit.QueueMsg{
						TaskId:       rabbit.GAME_START,
						GameServerID: gs.ID,
						Game:         gameDef,
					}

					err := SendRpcMessage("wrapper_"+gs.ID.String(), msg)
					if err != nil {
						logrus.Error(err.Error())
						break
					}

					model.Redis.Del("gs_restart_id_" + gs.ID.String())
				}

				_, err = model.Redis.Get("gs_start_id_" + gameServerId.String()).Int64()
				if err == nil {
					var gs model.GameServer
					if model.DB.Where("id = ?", gameServerId).First(&gs).RecordNotFound() {
						break
					}

					var gameDef game.Game
					err = json.Unmarshal([]byte(gs.GameJson), &gameDef)

					msg := rabbit.QueueMsg{
						TaskId:       rabbit.GAME_START,
						GameServerID: gs.ID,
						Game:         gameDef,
					}

					err := SendRpcMessage("wrapper_"+gs.ID.String(), msg)
					if err != nil {
						logrus.Error(err.Error())
						break
					}

					model.Redis.Del("gs_start_id_" + gs.ID.String())
				}

			case rabbit.WRAPPER_EXITED:
				logrus.Info("WRAPPER_EXITED")
				gameServerId := msgBody.GameServerID

				val, err := model.Redis.Get("gs_restart_id_" + gameServerId.String()).Int()
				if err != redis.Nil && val != -1 {
					model.Redis.Set("gs_restart_id_"+gameServerId.String(), 2, 24*time.Hour)

					var gs model.GameServer
					if model.DB.Where("id = ?", gameServerId).First(&gs).RecordNotFound() {
						break
					}

					var host model.Host
					if model.DB.Where("id = ?", gs.HostId).First(&host).RecordNotFound() {
						break
					}

					msg := rabbit.QueueMsg{
						TaskId:       rabbit.WRAPPER_START,
						GameServerID: gs.ID,
					}

					err := SendRpcMessage("agent_"+host.Token, msg)
					if err != nil {
						logrus.Error(err.Error())
						break
					}

					model.Redis.Set("gs_restart_id_"+gameServerId.String(), 3, 24*time.Hour)
				}

			case rabbit.WRAPPER_METRICS_FREQUENCY:
				gameServerId := msgBody.GameServerID

				var gs model.GameServer
				if model.DB.Where("id = ?", gameServerId).First(&gs).RecordNotFound() {
					break
				}

				msg := rabbit.QueueMsg{
					TaskId:          rabbit.WRAPPER_METRICS_FREQUENCY,
					MetricFrequency: gs.MetricFrequency,
				}

				err := SendRpcMessage("wrapper_"+gameServerId.String(), msg)
				if err != nil {
					logrus.Error(err.Error())
				}
			case rabbit.WRAPPER_METRICS:
				gameServerId := msgBody.GameServerID

				metric := &model.MetricGameServer{
					GameServerId: gameServerId,
					CpuUsage:     msgBody.CpuUsage,
					RamUsage:     msgBody.RamUsage,
				}
				model.DB.Save(metric)
			case rabbit.WRAPPER_HEARTBEAT:
				gameServerId := msgBody.GameServerID
				model.Redis.Set("wrapper_heartbeat_id_"+gameServerId.String(), msgBody.Timestamp, 24*time.Hour)
			}
		}
	}()
}

func ListenAgentData() {
	queue, err := channel.QueueDeclare(
		"agent_data",
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	lib.FailOnError(err, "Failed to declare a queue")

	msgs, err := channel.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil)
	lib.FailOnError(err, "Failed to register a consumer")

	go func() {
		for msg := range msgs {
			var msgBody rabbit.QueueMsg
			err = json.Unmarshal(msg.Body, &msgBody)
			if err != nil {
				logrus.Warn(err)
			}

			switch msgBody.TaskId {
			case rabbit.AGENT_METRICS_FREQUENCY:
				agentToken := msgBody.AgentToken

				var host model.Host
				if model.DB.Where("token = ?", agentToken).First(&host).RecordNotFound() {
					break
				}

				msg := rabbit.QueueMsg{
					TaskId:          rabbit.AGENT_METRICS_FREQUENCY,
					MetricFrequency: host.MetricFrequency,
				}

				err := SendRpcMessage("agent_"+agentToken, msg)
				if err != nil {
					logrus.Error(err.Error())
				}
			case rabbit.AGENT_METRICS:
				agentToken := msgBody.AgentToken

				var host model.Host
				if model.DB.Where("token = ?", agentToken).First(&host).RecordNotFound() {
					break
				}

				metric := &model.MetricHost{
					HostId:    host.ID,
					CpuUsage:  msgBody.CpuUsage,
					RamFree:   msgBody.RamFree,
					RamTotal:  msgBody.RamTotal,
					DiskFree:  msgBody.DiskFree,
					DiskUsed:  msgBody.DiskUsed,
					DiskTotal: msgBody.DiskTotal,
					User:      msgBody.User,
					System:    msgBody.System,
					Idle:      msgBody.Idle,
					Nice:      msgBody.Nice,
					Iowait:    msgBody.Iowait,
					Irq:       msgBody.Irq,
					Softirq:   msgBody.Softirq,
					Steal:     msgBody.Steal,
					Guest:     msgBody.Guest,
					GuestNice: msgBody.GuestNice,
				}
				model.DB.Save(metric)
			case rabbit.AGENT_OS:
				agentToken := msgBody.AgentToken

				var host model.Host
				if model.DB.Where("token = ?", agentToken).First(&host).RecordNotFound() {
					break
				}

				host.OS = msgBody.OS
				host.Platform = msgBody.Platform
				host.PlatformFamily = msgBody.PlatformFamily
				host.PlatformVersion = msgBody.PlatformVersion
				host.KernelVersion = msgBody.KernelVersion
				host.KernelArch = msgBody.KernelArch

				model.DB.Save(&host)

			case rabbit.AGENT_HEARTBEAT:
				agentToken := msgBody.AgentToken
				model.Redis.Set("agent_heartbeat_token_"+agentToken, msgBody.Timestamp, 24*time.Hour)
			}
		}
	}()
}
