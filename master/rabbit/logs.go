package rabbit

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/lib"
	"gitlab.com/systemz/aimpanel2/lib/rabbit"
	"gitlab.com/systemz/aimpanel2/master/db"
	"gitlab.com/systemz/aimpanel2/master/model"
	"gitlab.com/systemz/aimpanel2/master/redis"
	"time"
)

func ListenWrapperLogsQueue() {
	msgs, err := channel.Consume(
		"wrapper_logs",
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

			logrus.Info(msgBody)

			switch msgBody.TaskId {
			case rabbit.SERVER_LOG:
				var gsLog model.GameServerLog
				gsLog.GameServerID = msgBody.GameServerID

				if len(msgBody.Stdout) > 0 {
					gsLog.Log = msgBody.Stdout
					gsLog.Type = model.STDOUT
				}

				if len(msgBody.Stderr) > 0 {
					gsLog.Log = msgBody.Stderr
					gsLog.Type = model.STDERR
				}

				err = db.DB.Save(&gsLog).Error
				if err != nil {
					logrus.Warn(err)
				}

				logrus.Printf("Received a message: %s", gsLog)
			case rabbit.WRAPPER_STARTED:
				logrus.Info("WRAPPER_STARTED")
				gameServerId := msgBody.GameServerID
				_, err := redis.Redis.Get("gs_restart_id_" + gameServerId.String()).Int()
				if err == nil {
					var gs model.GameServer
					if db.DB.Where("id = ?", gameServerId).First(&gs).RecordNotFound() {
						break
					}

					var startCommand model.GameCommand
					if db.DB.Where("game_id = ? and type = ?", gs.GameId, "start").
						First(&startCommand).RecordNotFound() {
						break
					}

					msg := rabbit.QueueMsg{
						TaskId:           rabbit.GAME_START,
						GameServerID:     gs.ID,
						GameStartCommand: &startCommand,
					}

					err := SendRpcMessage("wrapper_"+gs.ID.String(), msg)
					if err != nil {
						logrus.Error(err.Error())
						break
					}

					redis.Redis.Del("gs_restart_id_" + gs.ID.String())
				}

				_, err = redis.Redis.Get("gs_start_id_" + gameServerId.String()).Int()
				if err == nil {
					var gs model.GameServer
					if db.DB.Where("id = ?", gameServerId).First(&gs).RecordNotFound() {
						break
					}

					var startCommand model.GameCommand
					if db.DB.Where("game_id = ? and type = ?", gs.GameId, "start").
						First(&startCommand).RecordNotFound() {
						break
					}

					msg := rabbit.QueueMsg{
						TaskId:           rabbit.GAME_START,
						GameServerID:     gs.ID,
						GameStartCommand: &startCommand,
					}

					err := SendRpcMessage("wrapper_"+gs.ID.String(), msg)
					if err != nil {
						logrus.Error(err.Error())
						break
					}

					redis.Redis.Del("gs_start_id_" + gs.ID.String())
				}

			case rabbit.WRAPPER_EXITED:
				logrus.Info("WRAPPER_EXITED")
				gameServerId := msgBody.GameServerID

				_, err := redis.Redis.Get("gs_restart_id_" + gameServerId.String()).Int()
				if err == nil {
					redis.Redis.Set("gs_restart_id_"+gameServerId.String(), 2, 1*time.Hour)

					var gs model.GameServer
					if db.DB.Where("id = ?", gameServerId).First(&gs).RecordNotFound() {
						break
					}

					var host model.Host
					if db.DB.Where("id = ?", gs.HostId).First(&host).RecordNotFound() {
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

					redis.Redis.Set("gs_restart_id_"+gameServerId.String(), 3, 1*time.Hour)
				}
			}
		}
	}()
}
