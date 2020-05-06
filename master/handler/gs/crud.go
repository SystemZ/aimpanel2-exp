package gs

import (
	"encoding/json"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"gitlab.com/systemz/aimpanel2/lib"
	"gitlab.com/systemz/aimpanel2/lib/ecode"
	"gitlab.com/systemz/aimpanel2/lib/game"
	"gitlab.com/systemz/aimpanel2/lib/request"
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
	hostId := params["host_id"]
	oid, _ := primitive.ObjectIDFromHex(hostId)

	//Decode json
	data := &request.GameServerCreate{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		lib.ReturnError(w, http.StatusBadRequest, ecode.JsonDecode, nil)
		return
	}

	gameDef := game.Game{
		Id:      data.GameId,
		Version: data.GameVersion,
	}
	gameDef.SetDefaults()
	gameDefJson, _ := json.Marshal(gameDef)
	//gameServer.GameJson = string(gameDefJson)

	//Check if host exist
	host := model.GetHost(oid)
	if host == nil {
		lib.ReturnError(w, http.StatusBadRequest, ecode.HostNotFound, nil)
		return
	}

	gameServer := &model.GameServer{
		Name:            data.Name,
		GameId:          data.GameId,
		GameVersion:     data.GameVersion,
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

	//TODO: create array of permissions?
	user := context.Get(r, "user").(model.User)
	group := model.GetGroup("USER-" + user.ID.String())
	if group == nil {
		lib.ReturnError(w, http.StatusInternalServerError, ecode.GroupNotFound, nil)
		return
	}

	// FIXME handle errors
	model.CreatePermissionsForNewGameServer(group.ID, host.ID, gameServer.ID)

	lib.MustEncode(json.NewEncoder(w),
		response.ID{ID: gameServer.ID.String()})
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
	hostId := params["id"]
	oid, _ := primitive.ObjectIDFromHex(hostId)

	host := model.GetHost(oid)
	if host == nil {
		lib.ReturnError(w, http.StatusNoContent, ecode.HostNotFound, nil)
		return
	}

	gameServers := model.GetGameServersByHostId(host.ID)
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

	gameServers := model.GetUserGameServers(user.ID)
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
	serverId, _ := primitive.ObjectIDFromHex(params["server_id"])
	hostId, _ := primitive.ObjectIDFromHex(params["hostId"])
	gameServer := model.GetGameServerByGsIdAndHostId(serverId, hostId)
	lib.MustEncode(json.NewEncoder(w), response.GameServer{GameServer: *gameServer})
}

func ConsoleLog(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	gameServerId := params["server_id"]

	logs := model.GetLogsByGameServer(gameServerId, 20)
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
	gameServerId, _ := primitive.ObjectIDFromHex(params["server_id"])
	err := gameserver.Remove(gameServerId)
	if err != nil {
		lib.ReturnError(w, http.StatusInternalServerError, ecode.GsRemove, err)
		return
	}

	lib.MustEncode(json.NewEncoder(w), response.JsonSuccess{Message: "Removing game server"})
}
