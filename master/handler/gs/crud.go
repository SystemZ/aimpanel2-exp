package gs

import (
	"encoding/json"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"gitlab.com/systemz/aimpanel2/lib"
	"gitlab.com/systemz/aimpanel2/lib/ecode"
	"gitlab.com/systemz/aimpanel2/lib/game"
	"gitlab.com/systemz/aimpanel2/lib/request"
	"gitlab.com/systemz/aimpanel2/master/handler"
	"gitlab.com/systemz/aimpanel2/master/model"
	"gitlab.com/systemz/aimpanel2/master/response"
	"gitlab.com/systemz/aimpanel2/master/service/gameserver"
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
// @Failure 400 {object} handler.JsonError
// @Security ApiKey
func Create(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	hostId := params["host_id"]

	//Decode json
	data := &request.GameServerCreate{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		lib.MustEncode(json.NewEncoder(w),
			handler.JsonError{ErrorCode: ecode.JsonDecode})
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
	host := model.GetHost(hostId)
	if host == nil {
		lib.MustEncode(json.NewEncoder(w),
			handler.JsonError{ErrorCode: ecode.HostNotFound})
		return
	}

	gameServer := &model.GameServer{
		Base: model.Base{
			DocType: "game_server",
		},
		Name:            data.Name,
		GameId:          data.GameId,
		GameVersion:     data.GameVersion,
		HostId:          host.ID,
		MetricFrequency: 30,
		StopTimeout:     30,
		GameJson:        string(gameDefJson),
	}

	//Save game server to db
	err = gameServer.Put(&gameServer)
	if err != nil {
		lib.MustEncode(json.NewEncoder(w),
			handler.JsonError{ErrorCode: ecode.DbError})
		return
	}

	//TODO: create array of permissions?
	user := context.Get(r, "user").(model.User)
	group := model.GetGroup("USER-" + user.ID)
	if group == nil {
		lib.MustEncode(json.NewEncoder(w),
			handler.JsonError{ErrorCode: ecode.GroupNotFound})
		return
	}

	// FIXME handle errors
	model.CreatePermissionsForNewGameServer(group.ID, host.ID, gameServer.ID)

	lib.MustEncode(json.NewEncoder(w),
		response.ID{ID: gameServer.ID})
}

// @Router /host/{id}/server [get]
// @Summary Game server list by Host ID
// @Tags Game Server
// @Description List game servers linked to selected host
// @Accept json
// @Produce json
// @Success 200 {object} response.GameServerList
// @Failure 400 {object} handler.JsonError
// @Security ApiKey
func ListByHostId(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	hostId := params["id"]

	host := model.GetHost(hostId)
	if host == nil {
		lib.MustEncode(json.NewEncoder(w),
			handler.JsonError{ErrorCode: ecode.HostNotFound})
		return
	}

	gameServers := model.GetGameServersByHostId(host.ID)
	if gameServers == nil {
		lib.MustEncode(json.NewEncoder(w),
			handler.JsonError{ErrorCode: ecode.GsNotFound})
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
// @Failure 400 {object} handler.JsonError
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
// @Failure 400 {object} handler.JsonError
// @Security ApiKey
func Get(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	gameServer := model.GetGameServerByGsIdAndHostId(params["server_id"], params["host_id"])
	lib.MustEncode(json.NewEncoder(w), response.GameServer{GameServer: *gameServer})
}

func ConsoleLog(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	gameServerId := params["server_id"]

	logs := model.GetLogsByGameServer(gameServerId, 20)
	if logs == nil {
		lib.MustEncode(json.NewEncoder(w),
			handler.JsonError{ErrorCode: ecode.GsNoLogs})
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
// @Success 200 {object} handler.JsonSuccess
// @Failure 400 {object} handler.JsonError
// @Security ApiKey
func Remove(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	gameServerId := params["server_id"]
	err := gameserver.Remove(gameServerId)
	if err != nil {
		lib.MustEncode(json.NewEncoder(w),
			handler.JsonError{ErrorCode: ecode.GsRemove})
		return
	}

	lib.MustEncode(json.NewEncoder(w), handler.JsonSuccess{Message: "Removing game server"})
}
