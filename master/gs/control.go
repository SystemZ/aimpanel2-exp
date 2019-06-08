package gs

import (
	"gitlab.com/systemz/aimpanel2/lib"
	"gitlab.com/systemz/aimpanel2/lib/rabbit"
	"gitlab.com/systemz/aimpanel2/master/model"
	rabbitMaster "gitlab.com/systemz/aimpanel2/master/rabbit"
	"time"
)

func Start(hostId string, gsId string, user model.User) error {
	host := user.GetHost(model.DB, hostId)
	if host == nil {
		return &lib.Error{ErrorCode: 5001}
	}

	gameServer := host.GetGameServer(model.DB, gsId)
	if gameServer == nil {
		return &lib.Error{ErrorCode: 5002}
	}

	game := gameServer.GetGame(model.DB)
	if game == nil {
		return &lib.Error{ErrorCode: 5003}
	}

	startCommand := game.GetStartCommandByVersion(model.DB, gameServer.GameVersion)
	if startCommand == nil {
		return &lib.Error{ErrorCode: 5004}
	}

	model.Redis.Set("gs_start_id_"+gameServer.ID.String(), 0, 1*time.Hour)

	err := rabbitMaster.SendRpcMessage("agent_"+host.Token, rabbit.QueueMsg{
		TaskId:       rabbit.WRAPPER_START,
		Game:         game.Name,
		GameServerID: gameServer.ID,
	})
	if err != nil {
		return &lib.Error{ErrorCode: 5005}
	}

	model.Redis.Set("gs_start_id_"+gameServer.ID.String(), 1, 1*time.Hour)

	return nil
}

func Stop(hostId string, gsId string, user model.User, stopType uint) error {
	host := user.GetHost(model.DB, hostId)
	if host == nil {
		return &lib.Error{ErrorCode: 5018}
	}

	gameServer := host.GetGameServer(model.DB, gsId)
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

func Install(hostId string, gsId string, user model.User) error {
	host := user.GetHost(model.DB, hostId)
	if host == nil {
		return &lib.Error{ErrorCode: 5006}
	}

	gameServer := host.GetGameServer(model.DB, gsId)
	if gameServer == nil {
		return &lib.Error{ErrorCode: 5007}
	}

	game := gameServer.GetGame(model.DB)
	if game == nil {
		return &lib.Error{ErrorCode: 5008}
	}

	gameFile := game.GetInstallFileByVersion(model.DB, gameServer.GameVersion)
	if gameFile == nil {
		return &lib.Error{ErrorCode: 5009}
	}

	installCommands := game.GetInstallCommandsByVersion(model.DB, gameServer.GameVersion)
	if installCommands == nil {
		return &lib.Error{ErrorCode: 5010}
	}

	err := rabbitMaster.SendRpcMessage("agent_"+host.Token, rabbit.QueueMsg{
		TaskId:       rabbit.GAME_INSTALL,
		Game:         game.Name,
		GameServerID: gameServer.ID,
		GameFile:     gameFile,
		GameCommands: installCommands,
	})
	if err != nil {
		return &lib.Error{ErrorCode: 5011}
	}

	return nil
}

func SendCommand(hostId string, gsId string, user model.User, command string) error {
	host := user.GetHost(model.DB, hostId)
	if host == nil {
		return &lib.Error{ErrorCode: 5027}
	}

	gameServer := host.GetGameServer(model.DB, gsId)
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

func Restart(hostId string, gsId string, user model.User, stopType uint) error {
	host := user.GetHost(model.DB, hostId)
	if host == nil {
		return &lib.Error{ErrorCode: 5013}
	}

	gameServer := host.GetGameServer(model.DB, gsId)
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
