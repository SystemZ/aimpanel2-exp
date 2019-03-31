package handler

import (
	"encoding/json"
	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	rabbit2 "gitlab.com/systemz/aimpanel2/lib/rabbit"
	"gitlab.com/systemz/aimpanel2/master/db"
	"gitlab.com/systemz/aimpanel2/master/model"
	"gitlab.com/systemz/aimpanel2/master/rabbit"
	"gitlab.com/systemz/aimpanel2/master/response"
	"net/http"
	"time"
)

func ListUserGameServersByHostId(w http.ResponseWriter, r *http.Request) {
	// swagger:route POST /hosts/{id}/servers GameServer list
	//
	// List user game servers by host id
	//
	//Consumes:
	//	- application/json
	//
	//Produces:
	//	- application/json
	//
	//Schemes: http, https
	//
	//Responses:
	//	default: jsonError
	//	200:

	userId := uuid.FromStringOrNil(r.Header.Get("uid"))
	params := mux.Vars(r)
	hostId := params["id"]

	var host model.Host
	if !db.DB.Where("id = ? and user_id = ?", hostId, userId).First(&host).RecordNotFound() {
		var gameServers []model.GameServer
		db.DB.Where("host_id = ?", hostId).Find(&gameServers)

		json.NewEncoder(w).Encode(gameServers)
	}
}

func ListUserGameServers(w http.ResponseWriter, r *http.Request) {
	userId := uuid.FromStringOrNil(r.Header.Get("uid"))

	var gameServers []model.GameServer
	db.DB.Table("game_servers").Select("game_servers.*").Joins(
		"LEFT JOIN hosts ON game_servers.host_id = hosts.id").Where(
		"hosts.user_id = ?", userId).Find(&gameServers)

	json.NewEncoder(w).Encode(&gameServers)
}

func CreateGameServer(w http.ResponseWriter, r *http.Request) {
	// swagger:route POST /hosts/{id}/servers GameServer create
	//
	// Creates new game server
	//
	//Consumes:
	//	- application/json
	//
	//Produces:
	//	- application/json
	//
	//Schemes: http, https
	//
	//Responses:
	//	default: jsonError
	//	200:
	gameServer := &model.GameServer{}

	err := json.NewDecoder(r.Body).Decode(gameServer)
	if err != nil {
		json.NewEncoder(w).Encode(response.JsonError{ErrorCode: 10, Message: "Invalid body."})
		return
	}

	params := mux.Vars(r)

	gameServer.HostId = uuid.FromStringOrNil(params["id"])

	db.DB.Save(gameServer)

	json.NewEncoder(w).Encode(gameServer)
}

func StartGameServer(w http.ResponseWriter, r *http.Request) {
	// swagger:route GET /hosts/{host_id}/server/{server_id}/start Wrapper start
	//
	// Creates new game server
	//
	//Consumes:
	//	- application/json
	//
	//Produces:
	//	- application/json
	//
	//Schemes: http, https
	//
	//Responses:
	//	default: jsonError
	//	200:

	params := mux.Vars(r)

	userId := uuid.FromStringOrNil(r.Header.Get("uid"))
	hostId := params["host_id"]
	gameServerId := params["server_id"]

	var host model.Host
	if !db.DB.Where("id = ? and user_id = ?", hostId, userId).First(&host).RecordNotFound() {
		var gameServer model.GameServer
		if !db.DB.Where("id = ? and host_id = ?", gameServerId, hostId).First(&gameServer).RecordNotFound() {
			var game model.Game
			if !db.DB.Where("id = ?", gameServer.GameId).First(&game).RecordNotFound() {
				var startCommand model.GameCommand
				if !db.DB.Where("game_id = ? and type = ?", gameServer.GameId, "start").First(&startCommand).RecordNotFound() {
					msg := rabbit2.QueueMsg{
						TaskId:           rabbit2.WRAPPER_START,
						Game:             game.Name,
						GameServerID:     gameServer.ID,
						GameStartCommand: startCommand,
					}
					rabbit.SendRpcMessage("agent_"+host.Token, msg)

					time.Sleep(5 * time.Second)

					msg = rabbit2.QueueMsg{
						TaskId:           rabbit2.GAME_START,
						GameServerID:     gameServer.ID,
						GameStartCommand: startCommand,
					}
					rabbit.SendRpcMessage("wrapper_"+gameServer.ID.String(), msg)

					json.NewEncoder(w).Encode(response.JsonSuccess{Message: "Started game server succesfully."})
				}
			}
		}
	}
}

func InstallGameServer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	userId := uuid.FromStringOrNil(r.Header.Get("uid"))
	hostId := params["host_id"]
	gameServerId := params["server_id"]

	var host model.Host
	//TODO: fix that ugly ifs :(
	if !db.DB.Where("id = ? and user_id = ?", hostId, userId).First(&host).RecordNotFound() {
		var gameServer model.GameServer
		if !db.DB.Where("id = ? and host_id = ?", gameServerId, hostId).First(&gameServer).RecordNotFound() {
			var game model.Game
			if !db.DB.Where("id = ?", gameServer.GameId).First(&game).RecordNotFound() {
				var gameFile model.GameFile
				if !db.DB.Where("game_id = ?", gameServer.GameId).First(&gameFile).RecordNotFound() {
					var installCommands []model.GameCommand
					if !db.DB.Where("game_id = ? and type = ?", gameServer.GameId, "install").Order("`order` asc").Find(&installCommands).RecordNotFound() {
						msg := rabbit2.QueueMsg{
							TaskId:       rabbit2.GAME_INSTALL,
							Game:         game.Name,
							GameServerID: gameServer.ID,
							GameFile:     gameFile,
							GameCommands: installCommands,
						}
						rabbit.SendRpcMessage("agent_"+host.Token, msg)

						json.NewEncoder(w).Encode(response.JsonSuccess{Message: "Installed game server succesfully."})
					}
				}
			}
		}
	}
}
