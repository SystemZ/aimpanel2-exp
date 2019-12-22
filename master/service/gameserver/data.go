package gameserver

import (
	"encoding/json"
	"errors"
	"github.com/alexandrevicenzi/go-sse"
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/lib/game"
	"gitlab.com/systemz/aimpanel2/lib/task"
	"gitlab.com/systemz/aimpanel2/master/events"
	"gitlab.com/systemz/aimpanel2/master/model"
	"strconv"
)

func HostData(hostToken string) error {
	host := model.GetHostByToken(model.DB, hostToken)
	if host == nil {
		return errors.New("error when getting host from db")
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
	}

	return nil
}
