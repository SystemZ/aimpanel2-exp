package gs

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/lib"
	"gitlab.com/systemz/aimpanel2/lib/game"
	"gitlab.com/systemz/aimpanel2/lib/rabbit"
	"gitlab.com/systemz/aimpanel2/master/model"
	rabbitMaster "gitlab.com/systemz/aimpanel2/master/rabbit"
	"time"
)

func Start(gsId string) error {
	gameServer := model.GetGameServer(model.DB, gsId)
	if gameServer == nil {
		return &lib.Error{ErrorCode: 5002}
	}

	hostToken := model.GetHostToken(model.DB, gameServer.HostId.String())
	if hostToken == "" {
		return &lib.Error{ErrorCode: 5003}
	}

	model.Redis.Set("gs_start_id_"+gameServer.ID.String(), 0, 1*time.Hour)

	err := rabbitMaster.SendRpcMessage("agent_"+hostToken, rabbit.QueueMsg{
		TaskId:       rabbit.WRAPPER_START,
		GameServerID: gameServer.ID,
	})

	if err != nil {
		return &lib.Error{ErrorCode: 5005}
	}

	model.Redis.Set("gs_start_id_"+gameServer.ID.String(), 1, 1*time.Hour)

	return nil
}

func Stop(gsId string, stopType uint) error {
	gameServer := model.GetGameServer(model.DB, gsId)
	if gameServer == nil {
		return &lib.Error{ErrorCode: 5019}
	}

	msg := rabbit.QueueMsg{
		GameServerID: gameServer.ID,
	}
	if stopType == 1 {
		msg.TaskId = rabbit.GAME_STOP_SIGKILL
	} else if stopType == 2 {
		msg.TaskId = rabbit.GAME_STOP_SIGTERM
	}

	err := rabbitMaster.SendRpcMessage("wrapper_"+gameServer.ID.String(), msg)
	if err != nil {
		return &lib.Error{ErrorCode: 5020}
	}

	return nil
}

func Install(gsId string) error {
	gameServer := model.GetGameServer(model.DB, gsId)
	if gameServer == nil {
		return &lib.Error{ErrorCode: 5007}
	}

	hostToken := model.GetHostToken(model.DB, gameServer.HostId.String())
	if hostToken == "" {
		return &lib.Error{ErrorCode: 5003}
	}

	gameFile := model.GetGameFileByGameIdAndVersion(model.DB, gameServer.GameId, gameServer.GameVersion)
	if gameFile == nil {
		return &lib.Error{ErrorCode: 5009}
	}

	var g game.Game
	err := json.Unmarshal([]byte(gameServer.GameJson), &g)
	if err != nil {
		logrus.Error(err)
	}
	g.DownloadUrl = gameFile.DownloadUrl

	err = rabbitMaster.SendRpcMessage("agent_"+hostToken, rabbit.QueueMsg{
		TaskId:       rabbit.GAME_INSTALL,
		GameServerID: gameServer.ID,
		Game:         g,
	})
	if err != nil {
		return &lib.Error{ErrorCode: 5011}
	}

	return nil
}

func SendCommand(gsId string, command string) error {
	gameServer := model.GetGameServer(model.DB, gsId)
	if gameServer == nil {
		return &lib.Error{ErrorCode: 5028}
	}

	err := rabbitMaster.SendRpcMessage("wrapper_"+gameServer.ID.String(), rabbit.QueueMsg{
		TaskId: rabbit.GAME_COMMAND,
		Body:   command,
	})
	if err != nil {
		return &lib.Error{ErrorCode: 5029}
	}

	return nil
}

func Restart(gsId string, stopType uint) error {
	gameServer := model.GetGameServer(model.DB, gsId)
	if gameServer == nil {
		return &lib.Error{ErrorCode: 5014}
	}

	model.Redis.Set("gs_restart_id_"+gameServer.ID.String(), 0, 24*time.Hour)

	msg := rabbit.QueueMsg{
		GameServerID: gameServer.ID,
	}

	if stopType == 1 {
		msg.TaskId = rabbit.GAME_STOP_SIGKILL
	} else if stopType == 2 {
		msg.TaskId = rabbit.GAME_STOP_SIGTERM
	}

	err := rabbitMaster.SendRpcMessage("wrapper_"+gameServer.ID.String(), msg)
	if err != nil {
		return &lib.Error{ErrorCode: 5015}
	}

	model.Redis.Set("gs_restart_id_"+gameServer.ID.String(), 1, 24*time.Hour)

	go func() {
		<-time.After(time.Duration(gameServer.StopTimeout) * time.Second)

		val, err := model.Redis.Get("gs_restart_id_" + gameServer.ID.String()).Int64()
		if err != nil {
			return
		}

		if val == 1 {
			model.Redis.Set("gs_restart_id_"+gameServer.ID.String(), -1, 24*time.Hour)
		}
	}()

	return nil
}

func Remove(gsId string) error {
	gameServer := model.GetGameServer(model.DB, gsId)
	if gameServer == nil {
		return &lib.Error{ErrorCode: 5014}
	}

	hostToken := model.GetHostToken(model.DB, gameServer.HostId.String())
	if hostToken == "" {
		return &lib.Error{ErrorCode: 5003}
	}

	if gameServer.State == 1 {
		msg := rabbit.QueueMsg{
			GameServerID: gameServer.ID,
			TaskId:       rabbit.GAME_STOP_SIGKILL,
		}

		err := rabbitMaster.SendRpcMessage("wrapper_"+gameServer.ID.String(), msg)
		if err != nil {
			return &lib.Error{ErrorCode: 5100}
		}
	}

	err := rabbitMaster.SendRpcMessage("agent_"+hostToken, rabbit.QueueMsg{
		TaskId:       rabbit.AGENT_REMOVE_GS,
		GameServerID: gameServer.ID,
	})
	if err != nil {
		return &lib.Error{ErrorCode: 5101}
	}

	model.DB.Where("endpoint LIKE ?", "/v1/host/"+gameServer.HostId.String()+"/server/"+gsId+"%").Delete(&model.Permission{})
	model.DB.Delete(&gameServer)

	return nil
}
