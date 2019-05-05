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

	msg := rabbit.QueueMsg{
		TaskId:       rabbit.WRAPPER_START,
		Game:         game.Name,
		GameServerID: gameServer.ID,
	}

	err := rabbitMaster.SendRpcMessage("agent_"+host.Token, msg)
	if err != nil {
		lib.MustEncode(json.NewEncoder(w), response.JsonError{ErrorCode: 5005, Message: "Could not start a wrapper."})
		return
	}

	//todo: change it to redis session

	time.Sleep(5 * time.Second)

	msg = rabbit.QueueMsg{
		TaskId:           rabbit.GAME_START,
		GameServerID:     gameServer.ID,
		GameStartCommand: startCommand,
	}
	err = rabbitMaster.SendRpcMessage("wrapper_"+gameServer.ID.String(), msg)
	if err != nil {
		lib.MustEncode(json.NewEncoder(w),
			response.JsonError{ErrorCode: 5006, Message: "Could not start a game."})
		return
	}

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
			response.JsonError{ErrorCode: 50012, Message: "Could not install game server."})
		return
	}

	lib.MustEncode(json.NewEncoder(w),
		response.JsonSuccess{Message: "Installed game server successfully."})
}

//func StopGameServer(w http.ResponseWriter, r *http.Request) {
//	params := mux.Vars(r)
//
//	userId := uuid.FromStringOrNil(r.Header.Get("uid"))
//	hostId := params["host_id"]
//	gameServerId := params["server_id"]
//
//	stopReq := &game_server.StopGameServerRequest{}
//
//	err := json.NewDecoder(r.Body).Decode(stopReq)
//	if err != nil {
//		json.NewEncoder(w).Encode(response.JsonError{ErrorCode: 10, Message: "Invalid body."})
//		return
//	}
//
//	var host model.Host
//	if db.DB.Where("id = ? and user_id = ?", hostId, userId).First(&host).RecordNotFound() {
//		json.NewEncoder(w).Encode(response.JsonError{ErrorCode: 20, Message: "Could not find host."})
//		return
//	}
//
//	var gameServer model.GameServer
//	if db.DB.Where("id = ? and host_id = ?", gameServerId, hostId).First(&gameServer).RecordNotFound() {
//		json.NewEncoder(w).Encode(response.JsonError{ErrorCode: 21, Message: "Could not find game server."})
//		return
//	}
//
//	if stopReq.Type == 1 {
//		//sigkill
//		msg := rabbit2.QueueMsg{
//			TaskId:       rabbit2.GAME_STOP_SIGKILL,
//			GameServerID: gameServer.ID,
//		}
//		rabbit.SendRpcMessage("wrapper_"+gameServer.ID.String(), msg)
//	} else if stopReq.Type == 2 {
//		//sigterm
//		msg := rabbit2.QueueMsg{
//			TaskId:       rabbit2.GAME_STOP_SIGTERM,
//			GameServerID: gameServer.ID,
//		}
//		rabbit.SendRpcMessage("wrapper_"+gameServer.ID.String(), msg)
//	}
//
//	json.NewEncoder(w).Encode(response.JsonSuccess{Message: "Stopped game server succesfully."})
//
//}
