package handler

import (
	"encoding/json"
	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"gitlab.com/systemz/aimpanel2/master/db"
	"gitlab.com/systemz/aimpanel2/master/model"
	"gitlab.com/systemz/aimpanel2/master/response"
	"net/http"
)

func ListUserGameServersByHostId(w http.ResponseWriter, r *http.Request) {
	// swagger:route POST /hosts/{id}/game_servers GameServer list
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
	var gameServers []model.GameServer

	userId := uuid.FromStringOrNil(r.Header.Get("uid"))
	params := mux.Vars(r)

	db.DB.Where("user_id = ? AND host_id = ?", userId, params["id"]).Find(&gameServers)

	json.NewEncoder(w).Encode(gameServers)
}

func CreateGameServer(w http.ResponseWriter, r *http.Request) {
	// swagger:route POST /hosts/{id}/game_servers GameServer create
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
