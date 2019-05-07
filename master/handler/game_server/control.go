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
			response.JsonError{ErrorCode: 5001})
		return
	}

	gameServer := host.GetGameServer(db.DB, gameServerId)
	if gameServer == nil {
		lib.MustEncode(json.NewEncoder(w),
			response.JsonError{ErrorCode: 5002})
		return
	}

	game := gameServer.GetGame(db.DB)
	if game == nil {
		lib.MustEncode(json.NewEncoder(w),
			response.JsonError{ErrorCode: 5003})
		return
	}

	startCommand := game.GetStartCommandByVersion(db.DB, gameServer.GameVersion)
	if startCommand == nil {
		lib.MustEncode(json.NewEncoder(w),
			response.JsonError{ErrorCode: 5004})
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
			response.JsonError{ErrorCode: 5005})
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
			response.JsonError{ErrorCode: 5006})
		return
	}

	gameServer := host.GetGameServer(db.DB, gameServerId)
	if gameServer == nil {
		lib.MustEncode(json.NewEncoder(w),
			response.JsonError{ErrorCode: 5007})
		return
	}

	game := gameServer.GetGame(db.DB)
	if game == nil {
		lib.MustEncode(json.NewEncoder(w),
			response.JsonError{ErrorCode: 5008})
		return
	}

	gameFile := game.GetInstallFileByVersion(db.DB, gameServer.GameVersion)
	if gameFile == nil {
		lib.MustEncode(json.NewEncoder(w),
			response.JsonError{ErrorCode: 5009})
		return
	}

	installCommands := game.GetInstallCommandsByVersion(db.DB, gameServer.GameVersion)
	if installCommands == nil {
		lib.MustEncode(json.NewEncoder(w),
			response.JsonError{ErrorCode: 5010})
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
			response.JsonError{ErrorCode: 5011})
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
			response.JsonError{ErrorCode: 5012})
		return
	}

	user := context.Get(r, "user").(model.User)
	host := user.GetHost(db.DB, hostId)
	if host == nil {
		lib.MustEncode(json.NewEncoder(w),
			response.JsonError{ErrorCode: 5013})
		return
	}

	gameServer := host.GetGameServer(db.DB, gameServerId)
	if gameServer == nil {
		lib.MustEncode(json.NewEncoder(w),
			response.JsonError{ErrorCode: 5014})
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
				response.JsonError{ErrorCode: 5015})
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
				response.JsonError{ErrorCode: 5016})
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
			response.JsonError{ErrorCode: 5017})
		return
	}

	user := context.Get(r, "user").(model.User)
	host := user.GetHost(db.DB, hostId)
	if host == nil {
		lib.MustEncode(json.NewEncoder(w),
			response.JsonError{ErrorCode: 5018})
		return
	}

	gameServer := host.GetGameServer(db.DB, gameServerId)
	if gameServer == nil {
		lib.MustEncode(json.NewEncoder(w),
			response.JsonError{ErrorCode: 5019})
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
				response.JsonError{ErrorCode: 5020})
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
				response.JsonError{ErrorCode: 5021})
			return
		}
	}

	lib.MustEncode(json.NewEncoder(w), response.JsonSuccess{Message: "Stopping the game server."})
}
