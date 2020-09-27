package gs

import (
	"encoding/json"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"gitlab.com/systemz/aimpanel2/lib"
	"gitlab.com/systemz/aimpanel2/lib/ecode"
	"gitlab.com/systemz/aimpanel2/lib/game"
	"gitlab.com/systemz/aimpanel2/lib/request"
	"gitlab.com/systemz/aimpanel2/lib/task"
	"gitlab.com/systemz/aimpanel2/master/model"
	"gitlab.com/systemz/aimpanel2/master/response"
	"gitlab.com/systemz/aimpanel2/master/service/gameserver"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

// @Router /host/{id}/server [post]
// @Summary Create
// @Tags Game Server
// @Description Create new game server linked to the selected host
// @Accept json
// @Produce json
// @Param host body request.GameServerCreate true " "
// @Success 200 {object} response.ID
// @Failure 400 {object} response.JsonError
// @Security ApiKey
func Create(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	hostId := params["hostId"]
	oid, err := primitive.ObjectIDFromHex(hostId)
	if err != nil {
		lib.ReturnError(w, http.StatusBadRequest, ecode.OidError, nil)
		return
	}

	//Decode json
	data := &request.GameServerCreate{}
	err = json.NewDecoder(r.Body).Decode(data)
	if err != nil {
		lib.ReturnError(w, http.StatusBadRequest, ecode.JsonDecode, nil)
		return
	}

	gameDef := game.Game{
		Id:      *data.GameId,
		Version: *data.GameVersion,
	}
	gameDef.SetDefaults()
	gameDefJson, _ := json.Marshal(gameDef)
	//gameServer.GameJson = string(gameDefJson)

	//Check if host exist
	host, err := model.GetHostById(oid)
	if err != nil {
		lib.ReturnError(w, http.StatusBadRequest, ecode.DbError, nil)
		return
	}

	if host == nil {
		lib.ReturnError(w, http.StatusBadRequest, ecode.HostNotFound, nil)
		return
	}

	gameServer := &model.GameServer{
		Name:            *data.Name,
		GameId:          *data.GameId,
		GameVersion:     *data.GameVersion,
		HostId:          host.ID,
		MetricFrequency: 30,
		StopTimeout:     30,
		GameJson:        string(gameDefJson),
	}

	//Save game server to db
	err = model.Put(gameServer)
	if err != nil {
		lib.ReturnError(w, http.StatusInternalServerError, ecode.DbError, err)
		return
	}
	/*
		user := context.Get(r, "user").(model.User)
		// FIXME handle errors
		err = model.CreatePermissionsForNewGameServer(user.ID, host.ID, gameServer.ID)
		if err != nil {
		lib.ReturnError(w, http.StatusInternalServerError, ecode.DbSave, nil)
		return
		}
	*/

	lib.MustEncode(json.NewEncoder(w),
		response.ID{ID: gameServer.ID.Hex()})
}

// @Router /host/{id}/server [get]
// @Summary Game server list by Host ID
// @Tags Game Server
// @Description List game servers linked to selected host
// @Accept json
// @Produce json
// @Success 200 {object} response.GameServerList
// @Failure 400 {object} response.JsonError
// @Security ApiKey
func ListByHostId(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	hostId := params["hostId"]
	oid, err := primitive.ObjectIDFromHex(hostId)
	if err != nil {
		lib.ReturnError(w, http.StatusBadRequest, ecode.OidError, err)
		return
	}

	host, err := model.GetHostById(oid)
	if err != nil {
		lib.ReturnError(w, http.StatusInternalServerError, ecode.DbError, err)
		return
	}

	if host == nil {
		lib.ReturnError(w, http.StatusNoContent, ecode.HostNotFound, nil)
		return
	}

	gameServers, err := model.GetGameServersByHostId(host.ID)
	if err != nil {
		lib.ReturnError(w, http.StatusInternalServerError, ecode.DbError, err)
		return
	}

	if gameServers == nil {
		lib.ReturnError(w, http.StatusNoContent, ecode.GsNotFound, nil)
		return
	}

	lib.MustEncode(json.NewEncoder(w),
		response.GameServerList{GameServers: *gameServers})
}

// @Router /host/my/server [get]
// @Summary User game servers
// @Tags Game Server
// @Description List game servers linked to the current signed-in account
// @Accept json
// @Produce json
// @Success 200 {object} response.GameServerList
// @Failure 400 {object} response.JsonError
// @Security ApiKey
func ListByUser(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user").(model.User)

	gameServers, err := model.GetUserGameServers(user)
	if err != nil {
		lib.ReturnError(w, http.StatusInternalServerError, ecode.DbError, err)
		return
	}

	if gameServers == nil {
		lib.MustEncode(json.NewEncoder(w), response.GameServerList{GameServers: nil})
		return
	}

	lib.MustEncode(json.NewEncoder(w), response.GameServerList{GameServers: *gameServers})
}

// @Router /host/{host_id}/server/{server_id} [get]
// @Summary Details
// @Tags Game Server
// @Description Get details about Game server with selected ID linked to the current signed-in account
// @Accept json
// @Produce json
// @Param host_id path string true "Host ID"
// @Param server_id path string true "Game Server ID"
// @Success 200 {object} response.GameServer
// @Failure 400 {object} response.JsonError
// @Security ApiKey
func Get(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	serverId, err := primitive.ObjectIDFromHex(params["gsId"])
	if err != nil {
		lib.ReturnError(w, http.StatusBadRequest, ecode.OidError, err)
		return
	}
	hostId, err := primitive.ObjectIDFromHex(params["hostId"])
	if err != nil {
		lib.ReturnError(w, http.StatusBadRequest, ecode.OidError, err)
		return
	}

	gameServer, err := model.GetGameServerByGsIdAndHostId(serverId, hostId)
	if err != nil {
		lib.ReturnError(w, http.StatusInternalServerError, ecode.DbError, err)
		return
	}

	// prevent null port list
	if gameServer.Ports == nil {
		emptyPorts := make([]model.GamePort, 0)
		gameServer.Ports = &emptyPorts
	}

	lib.MustEncode(json.NewEncoder(w), response.GameServer{GameServer: *gameServer})
}

// @Router /host/{hostId}/server/{gsId} [PUT]
// @Summary Edit
// @Tags Game Server
// @Description Edit Game server with selected ID
// @Accept json
// @Produce json
// @Param hostId path string true "Host ID"
// @Param gsId path string true "Game Server ID"
// @Param host body request.GameServerCreate true " "
// @Success 200 {object} response.GameServer
// @Failure 400 {object} response.JsonError
// @Security ApiKey
func Edit(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	serverId, err := primitive.ObjectIDFromHex(params["gsId"])
	if err != nil {
		lib.ReturnError(w, http.StatusBadRequest, ecode.OidError, err)
		return
	}
	hostId, err := primitive.ObjectIDFromHex(params["hostId"])
	if err != nil {
		lib.ReturnError(w, http.StatusBadRequest, ecode.OidError, err)
		return
	}
	gameServer, err := model.GetGameServerByGsIdAndHostId(serverId, hostId)
	if err != nil {
		lib.ReturnError(w, http.StatusInternalServerError, ecode.DbError, err)
		return
	}

	//Decode json
	data := &request.GameServerCreate{}
	err = json.NewDecoder(r.Body).Decode(data)
	if err != nil {
		lib.ReturnError(w, http.StatusBadRequest, ecode.JsonDecode, nil)
		return
	}

	// modify ports
	if data.Ports != nil {
		if data.SerializePorts() != gameServer.SerializePorts() {
			// FIXME validate
			// FIXME log this action
			// FIXME maybe this can be done in less code
			var newPortList []model.GamePort
			for _, port := range *data.Ports {
				newPortList = append(newPortList, model.GamePort{
					Protocol:      port.Protocol,
					Host:          port.Host,
					PortHost:      port.PortHost,
					PortContainer: port.PortContainer,
				})
			}
			gameServer.Ports = &newPortList
			model.Update(gameServer)
		}
	}

	// modify custom cmd to start server
	if data.CustomCmdStart != nil {
		// FIXME validate custom CMD
		// FIXME move to service
		user := context.Get(r, "user").(model.User)
		err = model.SaveAction(
			task.Message{
				TaskId:       task.GS_CMD_START_CHANGE,
				GameServerID: gameServer.ID.Hex(),
			},
			user,
			hostId,
			*data.CustomCmdStart,
			gameServer.CustomCmdStart,
		)
		if err != nil {
			lib.ReturnError(w, http.StatusInternalServerError, ecode.DbSave, err)
			return
		}
		gameServer.CustomCmdStart = *data.CustomCmdStart
		model.Update(gameServer)
	}

	if data.Name != nil {
		user := context.Get(r, "user").(model.User)
		err = model.SaveAction(
			task.Message{
				TaskId:       task.GS_NAME_CHANGE,
				GameServerID: gameServer.ID.Hex(),
			},
			user,
			hostId,
			*data.Name,
			gameServer.Name,
		)
		if err != nil {
			lib.ReturnError(w, http.StatusInternalServerError, ecode.DbSave, err)
			return
		}

		gameServer.Name = *data.Name
		model.Update(gameServer)
	}

	if data.GameId != nil && data.GameVersion != nil {
		//&& gameServer.GameId != *data.GameId || gameServer.GameVersion != *data.GameVersion
		user := context.Get(r, "user").(model.User)
		err = model.SaveAction(
			task.Message{
				TaskId:       task.GS_GAME_CHANGE,
				GameServerID: gameServer.ID.Hex(),
			},
			user,
			hostId,
			game.GetGameNameById(gameServer.GameId)+" "+*data.GameVersion,
			game.GetGameNameById(*data.GameId)+" "+*data.GameVersion,
		)
		if err != nil {
			lib.ReturnError(w, http.StatusInternalServerError, ecode.DbSave, err)
			return
		}

		gameDef := game.Game{
			Id:      *data.GameId,
			Version: *data.GameVersion,
		}
		gameDef.SetDefaults()
		gameDefJson, _ := json.Marshal(gameDef)

		gameServer.GameId = *data.GameId
		gameServer.GameVersion = *data.GameVersion
		gameServer.GameJson = string(gameDefJson)

		model.Update(gameServer)
	}

	lib.MustEncode(json.NewEncoder(w), response.GameServer{GameServer: *gameServer})
}

func ConsoleLog(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	serverId, err := primitive.ObjectIDFromHex(params["gsId"])
	if err != nil {
		lib.ReturnError(w, http.StatusBadRequest, ecode.OidError, err)
		return
	}

	logs, err := model.GetLogsByGameServerId(serverId, 20)
	if err != nil {
		lib.ReturnError(w, http.StatusInternalServerError, ecode.DbError, err)
		return
	}

	if logs == nil {
		lib.ReturnError(w, http.StatusNoContent, ecode.GsNoLogs, nil)
		return
	}

	lib.MustEncode(json.NewEncoder(w), logs)
}

func PutLogs(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}

// @Router /host/{host_id}/server/{server_id} [delete]
// @Summary Remove
// @Tags Game Server
// @Description Removes game server
// @Accept json
// @Produce json
// @Param host_id path string true "Host ID"
// @Param server_id path string true "Game Server ID"
// @Success 200 {object} response.JsonSuccess
// @Failure 400 {object} response.JsonError
// @Security ApiKey
func Remove(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	gameServerId, _ := primitive.ObjectIDFromHex(params["gsId"])
	user := context.Get(r, "user").(model.User)
	err := gameserver.Remove(gameServerId, user)
	if err != nil {
		lib.ReturnError(w, http.StatusInternalServerError, ecode.GsRemove, err)
		return
	}

	lib.MustEncode(json.NewEncoder(w), response.JsonSuccess{Message: "Removing game server"})
}
