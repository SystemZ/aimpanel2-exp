package game_server

import (
	"encoding/json"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"gitlab.com/systemz/aimpanel2/lib"
	"gitlab.com/systemz/aimpanel2/lib/rabbit"
	"gitlab.com/systemz/aimpanel2/master/handler"
	"gitlab.com/systemz/aimpanel2/master/model"
	rabbitMaster "gitlab.com/systemz/aimpanel2/master/rabbit"
	"net/http"
	"time"
)

type GameServerStopReq struct {
	Type uint `json:"type"`
}

type GameServerSendCommandReq struct {
	Command string `json:"command"`
}

func Start(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	hostId := params["host_id"]
	gameServerId := params["server_id"]

	user := context.Get(r, "user").(model.User)
	host := user.GetHost(model.DB, hostId)
	if host == nil {
		lib.MustEncode(json.NewEncoder(w),
			handler.JsonError{ErrorCode: 5001})
		return
	}

	gameServer := host.GetGameServer(model.DB, gameServerId)
	if gameServer == nil {
		lib.MustEncode(json.NewEncoder(w),
			handler.JsonError{ErrorCode: 5002})
		return
	}

	game := gameServer.GetGame(model.DB)
	if game == nil {
		lib.MustEncode(json.NewEncoder(w),
			handler.JsonError{ErrorCode: 5003})
		return
	}

	startCommand := game.GetStartCommandByVersion(model.DB, gameServer.GameVersion)
	if startCommand == nil {
		lib.MustEncode(json.NewEncoder(w),
			handler.JsonError{ErrorCode: 5004})
		return
	}

	model.Redis.Set("gs_start_id_"+gameServer.ID.String(), 0, 1*time.Hour)

	msg := rabbit.QueueMsg{
		TaskId:       rabbit.WRAPPER_START,
		Game:         game.Name,
		GameServerID: gameServer.ID,
	}

	err := rabbitMaster.SendRpcMessage("agent_"+host.Token, msg)
	if err != nil {
		lib.MustEncode(json.NewEncoder(w),
			handler.JsonError{ErrorCode: 5005})
		return
	}

	model.Redis.Set("gs_start_id_"+gameServer.ID.String(), 1, 1*time.Hour)

	lib.MustEncode(json.NewEncoder(w),
		handler.JsonSuccess{Message: "Started game server succesfully."})
}

func Install(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	hostId := params["host_id"]
	gameServerId := params["server_id"]

	user := context.Get(r, "user").(model.User)
	host := user.GetHost(model.DB, hostId)
	if host == nil {
		lib.MustEncode(json.NewEncoder(w),
			handler.JsonError{ErrorCode: 5006})
		return
	}

	gameServer := host.GetGameServer(model.DB, gameServerId)
	if gameServer == nil {
		lib.MustEncode(json.NewEncoder(w),
			handler.JsonError{ErrorCode: 5007})
		return
	}

	game := gameServer.GetGame(model.DB)
	if game == nil {
		lib.MustEncode(json.NewEncoder(w),
			handler.JsonError{ErrorCode: 5008})
		return
	}

	gameFile := game.GetInstallFileByVersion(model.DB, gameServer.GameVersion)
	if gameFile == nil {
		lib.MustEncode(json.NewEncoder(w),
			handler.JsonError{ErrorCode: 5009})
		return
	}

	installCommands := game.GetInstallCommandsByVersion(model.DB, gameServer.GameVersion)
	if installCommands == nil {
		lib.MustEncode(json.NewEncoder(w),
			handler.JsonError{ErrorCode: 5010})
		return
	}

	msg := rabbit.QueueMsg{
		TaskId:       rabbit.GAME_INSTALL,
		Game:         game.Name,
		GameServerID: gameServer.ID,
		GameFile:     gameFile,
		GameCommands: installCommands,
	}

	err := rabbitMaster.SendRpcMessage("agent_"+host.Token, msg)
	if err != nil {
		lib.MustEncode(json.NewEncoder(w),
			handler.JsonError{ErrorCode: 5011})
		return
	}

	lib.MustEncode(json.NewEncoder(w),
		handler.JsonSuccess{Message: "Installed game server successfully."})
}

func Restart(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	hostId := params["host_id"]
	gameServerId := params["server_id"]

	stopReq := &GameServerStopReq{}

	err := json.NewDecoder(r.Body).Decode(stopReq)
	if err != nil {
		lib.MustEncode(json.NewEncoder(w),
			handler.JsonError{ErrorCode: 5012})
		return
	}

	user := context.Get(r, "user").(model.User)
	host := user.GetHost(model.DB, hostId)
	if host == nil {
		lib.MustEncode(json.NewEncoder(w),
			handler.JsonError{ErrorCode: 5013})
		return
	}

	gameServer := host.GetGameServer(model.DB, gameServerId)
	if gameServer == nil {
		lib.MustEncode(json.NewEncoder(w),
			handler.JsonError{ErrorCode: 5014})
		return
	}

	model.Redis.Set("gs_restart_id_"+gameServer.ID.String(), 0, 1*time.Hour)

	msg := rabbit.QueueMsg{
		GameServerID: gameServer.ID,
	}

	if stopReq.Type == 1 {
		msg.TaskId = rabbit.GAME_STOP_SIGKILL
	} else if stopReq.Type == 2 {
		msg.TaskId = rabbit.GAME_STOP_SIGTERM
	}

	err = rabbitMaster.SendRpcMessage("wrapper_"+gameServer.ID.String(), msg)
	if err != nil {
		lib.MustEncode(json.NewEncoder(w),
			handler.JsonError{ErrorCode: 5015})
		return
	}

	model.Redis.Set("gs_restart_id_"+gameServer.ID.String(), 1, 1*time.Hour)

	lib.MustEncode(json.NewEncoder(w), handler.JsonSuccess{Message: "Restarting the game server."})
}

func Stop(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	hostId := params["host_id"]
	gameServerId := params["server_id"]

	stopReq := &GameServerStopReq{}

	err := json.NewDecoder(r.Body).Decode(stopReq)
	if err != nil {
		lib.MustEncode(json.NewEncoder(w),
			handler.JsonError{ErrorCode: 5017})
		return
	}

	user := context.Get(r, "user").(model.User)
	host := user.GetHost(model.DB, hostId)
	if host == nil {
		lib.MustEncode(json.NewEncoder(w),
			handler.JsonError{ErrorCode: 5018})
		return
	}

	gameServer := host.GetGameServer(model.DB, gameServerId)
	if gameServer == nil {
		lib.MustEncode(json.NewEncoder(w),
			handler.JsonError{ErrorCode: 5019})
		return
	}

	msg := rabbit.QueueMsg{
		GameServerID: gameServer.ID,
	}

	if stopReq.Type == 1 {
		msg.TaskId = rabbit.GAME_STOP_SIGKILL
	} else if stopReq.Type == 2 {
		msg.TaskId = rabbit.GAME_STOP_SIGTERM
	}

	err = rabbitMaster.SendRpcMessage("wrapper_"+gameServer.ID.String(), msg)
	if err != nil {
		lib.MustEncode(json.NewEncoder(w),
			handler.JsonError{ErrorCode: 5020})
		return
	}

	lib.MustEncode(json.NewEncoder(w), handler.JsonSuccess{Message: "Stopping the game server."})
}

func SendCommand(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	hostId := params["host_id"]
	gameServerId := params["server_id"]

	cmdReq := &GameServerSendCommandReq{}

	err := json.NewDecoder(r.Body).Decode(cmdReq)
	if err != nil {
		lib.MustEncode(json.NewEncoder(w),
			handler.JsonError{ErrorCode: 5026})
		return
	}

	user := context.Get(r, "user").(model.User)
	host := user.GetHost(model.DB, hostId)
	if host == nil {
		lib.MustEncode(json.NewEncoder(w),
			handler.JsonError{ErrorCode: 5027})
		return
	}

	gameServer := host.GetGameServer(model.DB, gameServerId)
	if gameServer == nil {
		lib.MustEncode(json.NewEncoder(w),
			handler.JsonError{ErrorCode: 5028})
		return
	}

	msg := rabbit.QueueMsg{
		TaskId: rabbit.GAME_COMMAND,
		Body:   cmdReq.Command,
	}

	err = rabbitMaster.SendRpcMessage("wrapper_"+gameServer.ID.String(), msg)
	if err != nil {
		lib.MustEncode(json.NewEncoder(w),
			handler.JsonError{ErrorCode: 5029})
		return
	}

	lib.MustEncode(json.NewEncoder(w), handler.JsonSuccess{Message: "Sending command to server"})
}
