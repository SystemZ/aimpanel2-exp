package handler

import (
	"encoding/json"
	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	rabbit2 "gitlab.com/systemz/aimpanel2/lib/rabbit"
	"gitlab.com/systemz/aimpanel2/master/db"
	"gitlab.com/systemz/aimpanel2/master/model"
	"gitlab.com/systemz/aimpanel2/master/rabbit"
	"gitlab.com/systemz/aimpanel2/master/response"
	"net/http"
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
				msg := rabbit2.QueueMsg{
					TaskId:       rabbit2.WRAPPER_START,
					Game:         game.Name,
					GameServerID: gameServer.ID,
				}
				rabbit.SendRpcMessage("agent_"+host.Token, msg)

				json.NewEncoder(w).Encode(response.JsonSuccess{Message: "Started game server succesfully."})
			}
		}
	}
}

func InstallGameServer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	userId := uuid.FromStringOrNil(r.Header.Get("uid"))
	hostId := params["host_id"]
	gameServerId := params["server_id"]

	logrus.Info("1")
	var host model.Host
	if !db.DB.Where("id = ? and user_id = ?", hostId, userId).First(&host).RecordNotFound() {
		logrus.Info("2")
		var gameServer model.GameServer
		if !db.DB.Where("id = ? and host_id = ?", gameServerId, hostId).First(&gameServer).RecordNotFound() {
			logrus.Info("3")
			var game model.Game
			if !db.DB.Where("id = ?", gameServer.GameId).First(&game).RecordNotFound() {
				logrus.Info("4")
				msg := rabbit2.QueueMsg{
					TaskId:       rabbit2.GAME_INSTALL,
					Game:         game.Name,
					GameServerID: gameServer.ID,
				}
				rabbit.SendRpcMessage("agent_"+host.Token, msg)

				json.NewEncoder(w).Encode(response.JsonSuccess{Message: "Installed game server succesfully."})
			}
		}
	}
}
