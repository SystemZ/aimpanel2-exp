package game_server

import (
	"encoding/json"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"gitlab.com/systemz/aimpanel2/lib"
	"gitlab.com/systemz/aimpanel2/lib/rabbit"
	"gitlab.com/systemz/aimpanel2/master/db"
	"gitlab.com/systemz/aimpanel2/master/model"
	rabbitMaster "gitlab.com/systemz/aimpanel2/master/rabbit"
	"gitlab.com/systemz/aimpanel2/master/redis"
	"gitlab.com/systemz/aimpanel2/master/request/game_server"
	"gitlab.com/systemz/aimpanel2/master/response"
	"net/http"
	"time"
)

func Start(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	hostId := params["host_id"]
	gameServerId := params["server_id"]

	user := context.Get(r, "user").(model.User)
	host := user.GetHost(db.DB, hostId)
	if host == nil {
		lib.MustEncode(json.NewEncoder(w),
			response.JsonError{ErrorCode: 5001, Message: "Could not find a host."})
		return
	}

	gameServer := host.GetGameServer(db.DB, gameServerId)
	if gameServer == nil {
		lib.MustEncode(json.NewEncoder(w),
			response.JsonError{ErrorCode: 5002, Message: "Could not find a game server."})
		return
	}

	game := gameServer.GetGame(db.DB)
	if game == nil {
		lib.MustEncode(json.NewEncoder(w),
			response.JsonError{ErrorCode: 5003, Message: "Could not find a game."})
		return
	}

	startCommand := game.GetStartCommand(db.DB)
	if startCommand == nil {
		lib.MustEncode(json.NewEncoder(w),
			response.JsonError{ErrorCode: 5004, Message: "Could not find a start command."})
		return
	}

	redis.Redis.Set("gs_start_id_"+gameServer.ID.String(), 0, 1*time.Hour)

	msg := rabbit.QueueMsg{
		TaskId:       rabbit.WRAPPER_START,
		Game:         game.Name,
		GameServerID: gameServer.ID,
	}

	err := rabbitMaster.SendRpcMessage("agent_"+host.Token, msg)
	if err != nil {
		lib.MustEncode(json.NewEncoder(w),
			response.JsonError{ErrorCode: 5005, Message: "Could not start a wrapper."})
		return
	}

	redis.Redis.Set("gs_start_id_"+gameServer.ID.String(), 1, 1*time.Hour)

	lib.MustEncode(json.NewEncoder(w),
		response.JsonSuccess{Message: "Started game server succesfully."})
}

func Install(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	hostId := params["host_id"]
	gameServerId := params["server_id"]

	user := context.Get(r, "user").(model.User)
	host := user.GetHost(db.DB, hostId)
	if host == nil {
		lib.MustEncode(json.NewEncoder(w),
			response.JsonError{ErrorCode: 5007, Message: "Could not find a host."})
		return
	}

	gameServer := host.GetGameServer(db.DB, gameServerId)
	if gameServer == nil {
		lib.MustEncode(json.NewEncoder(w),
			response.JsonError{ErrorCode: 5008, Message: "Could not find a game server."})
		return
	}

	game := gameServer.GetGame(db.DB)
	if game == nil {
		lib.MustEncode(json.NewEncoder(w),
			response.JsonError{ErrorCode: 5009, Message: "Could not find a game."})
		return
	}

	gameFile := game.GetInstallFile(db.DB)
	if gameFile == nil {
		lib.MustEncode(json.NewEncoder(w),
			response.JsonError{ErrorCode: 5010, Message: "Could not find a install file."})
		return
	}

	installCommands := game.GetInstallCommands(db.DB)
	if installCommands == nil {
		lib.MustEncode(json.NewEncoder(w),
			response.JsonError{ErrorCode: 5011, Message: "Could not find a install commands."})
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
			response.JsonError{ErrorCode: 5012, Message: "Could not install game server."})
		return
	}

	lib.MustEncode(json.NewEncoder(w),
		response.JsonSuccess{Message: "Installed game server successfully."})
}

func Restart(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	hostId := params["host_id"]
	gameServerId := params["server_id"]

	stopReq := &game_server.StopGameServerRequest{}

	err := json.NewDecoder(r.Body).Decode(stopReq)
	if err != nil {
		lib.MustEncode(json.NewEncoder(w),
			response.JsonError{ErrorCode: 5013, Message: "Invalid body."})
		return
	}

	user := context.Get(r, "user").(model.User)
	host := user.GetHost(db.DB, hostId)
	if host == nil {
		lib.MustEncode(json.NewEncoder(w),
			response.JsonError{ErrorCode: 5014, Message: "Could not find a host."})
		return
	}

	gameServer := host.GetGameServer(db.DB, gameServerId)
	if gameServer == nil {
		lib.MustEncode(json.NewEncoder(w),
			response.JsonError{ErrorCode: 5015, Message: "Could not find a game server."})
		return
	}

	redis.Redis.Set("gs_restart_id_"+gameServer.ID.String(), 0, 1*time.Hour)

	if stopReq.Type == 1 {
		//sigkill
		msg := rabbit.QueueMsg{
			TaskId:       rabbit.GAME_STOP_SIGKILL,
			GameServerID: gameServer.ID,
		}

		err = rabbitMaster.SendRpcMessage("wrapper_"+gameServer.ID.String(), msg)
		if err != nil {
			lib.MustEncode(json.NewEncoder(w),
				response.JsonError{ErrorCode: 5016, Message: "Could not stop a game."})
			return
		}
	} else if stopReq.Type == 2 {
		//sigterm
		msg := rabbit.QueueMsg{
			TaskId:       rabbit.GAME_STOP_SIGTERM,
			GameServerID: gameServer.ID,
		}

		err = rabbitMaster.SendRpcMessage("wrapper_"+gameServer.ID.String(), msg)
		if err != nil {
			lib.MustEncode(json.NewEncoder(w),
				response.JsonError{ErrorCode: 5017, Message: "Could not stop a game."})
			return
		}
	}

	redis.Redis.Set("gs_restart_id_"+gameServer.ID.String(), 1, 1*time.Hour)

	lib.MustEncode(json.NewEncoder(w), response.JsonSuccess{Message: "Restarting the game server."})
}

func Stop(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	hostId := params["host_id"]
	gameServerId := params["server_id"]

	stopReq := &game_server.StopGameServerRequest{}

	err := json.NewDecoder(r.Body).Decode(stopReq)
	if err != nil {
		lib.MustEncode(json.NewEncoder(w),
			response.JsonError{ErrorCode: 5018, Message: "Invalid body."})
		return
	}

	user := context.Get(r, "user").(model.User)
	host := user.GetHost(db.DB, hostId)
	if host == nil {
		lib.MustEncode(json.NewEncoder(w),
			response.JsonError{ErrorCode: 5019, Message: "Could not find a host."})
		return
	}

	gameServer := host.GetGameServer(db.DB, gameServerId)
	if gameServer == nil {
		lib.MustEncode(json.NewEncoder(w),
			response.JsonError{ErrorCode: 5020, Message: "Could not find a game server."})
		return
	}

	if stopReq.Type == 1 {
		//sigkill
		msg := rabbit.QueueMsg{
			TaskId:       rabbit.GAME_STOP_SIGKILL,
			GameServerID: gameServer.ID,
		}

		err = rabbitMaster.SendRpcMessage("wrapper_"+gameServer.ID.String(), msg)
		if err != nil {
			lib.MustEncode(json.NewEncoder(w),
				response.JsonError{ErrorCode: 5021, Message: "Could not stop a game."})
			return
		}
	} else if stopReq.Type == 2 {
		//sigterm
		msg := rabbit.QueueMsg{
			TaskId:       rabbit.GAME_STOP_SIGTERM,
			GameServerID: gameServer.ID,
		}

		err = rabbitMaster.SendRpcMessage("wrapper_"+gameServer.ID.String(), msg)
		if err != nil {
			lib.MustEncode(json.NewEncoder(w),
				response.JsonError{ErrorCode: 5022, Message: "Could not stop a game."})
			return
		}
	}

	lib.MustEncode(json.NewEncoder(w), response.JsonSuccess{Message: "Stopping the game server."})
}
