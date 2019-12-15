package game_server

import (
	"encoding/json"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"gitlab.com/systemz/aimpanel2/lib"
	"gitlab.com/systemz/aimpanel2/lib/game"
	"gitlab.com/systemz/aimpanel2/master/gs"
	"gitlab.com/systemz/aimpanel2/master/handler"
	"gitlab.com/systemz/aimpanel2/master/model"
	"net/http"
)

func Create(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	hostId := params["host_id"]

	//Decode json
	gameServer := &model.GameServer{}
	err := json.NewDecoder(r.Body).Decode(gameServer)
	if err != nil {
		lib.MustEncode(json.NewEncoder(w),
			handler.JsonError{ErrorCode: 5022})
		return
	}

	gameDef := game.Game{
		Id:      gameServer.GameId,
		Version: gameServer.GameVersion,
	}
	gameDef.SetDefaults()

	gameDefJson, _ := json.Marshal(gameDef)
	gameServer.GameJson = string(gameDefJson)

	//Check if host exist
	host := model.GetHost(model.DB, hostId)
	if host == nil {
		lib.MustEncode(json.NewEncoder(w),
			handler.JsonError{ErrorCode: 5023})
		return
	}

	gameServer.HostId = host.ID
	gameServer.MetricFrequency = 30
	gameServer.StopTimeout = 30

	//Save game server to db
	model.DB.Save(gameServer)

	//TODO: create array of permissions?
	user := context.Get(r, "user").(model.User)
	group := model.GetGroup(model.DB, "USER-"+user.ID.String())
	if group == nil {
		lib.MustEncode(json.NewEncoder(w),
			handler.JsonError{ErrorCode: 3002})
		return
	}

	model.DB.Save(&model.Permission{
		Name:     "Get game server",
		Verb:     lib.GetVerbByName("GET"),
		GroupId:  group.ID,
		Endpoint: "/v1/host/" + host.ID.String() + "/server/" + gameServer.ID.String(),
	})

	model.DB.Save(&model.Permission{
		Name:     "Delete game server",
		Verb:     lib.GetVerbByName("DELETE"),
		GroupId:  group.ID,
		Endpoint: "/v1/host/" + host.ID.String() + "/server/" + gameServer.ID.String(),
	})

	model.DB.Save(&model.Permission{
		Name:     "Install game server",
		Verb:     lib.GetVerbByName("PUT"),
		GroupId:  group.ID,
		Endpoint: "/v1/host/" + host.ID.String() + "/server/" + gameServer.ID.String() + "/install",
	})

	model.DB.Save(&model.Permission{
		Name:     "Start game server",
		Verb:     lib.GetVerbByName("PUT"),
		GroupId:  group.ID,
		Endpoint: "/v1/host/" + host.ID.String() + "/server/" + gameServer.ID.String() + "/start",
	})

	model.DB.Save(&model.Permission{
		Name:     "Restart game server",
		Verb:     lib.GetVerbByName("PUT"),
		GroupId:  group.ID,
		Endpoint: "/v1/host/" + host.ID.String() + "/server/" + gameServer.ID.String() + "/restart",
	})

	model.DB.Save(&model.Permission{
		Name:     "Stop game server",
		Verb:     lib.GetVerbByName("PUT"),
		GroupId:  group.ID,
		Endpoint: "/v1/host/" + host.ID.String() + "/server/" + gameServer.ID.String() + "/stop",
	})

	model.DB.Save(&model.Permission{
		Name:     "Send command to game server",
		Verb:     lib.GetVerbByName("PUT"),
		GroupId:  group.ID,
		Endpoint: "/v1/host/" + host.ID.String() + "/server/" + gameServer.ID.String() + "/command",
	})

	model.DB.Save(&model.Permission{
		Name:     "Get logs from game server",
		Verb:     lib.GetVerbByName("PUT"),
		GroupId:  group.ID,
		Endpoint: "/v1/host/" + host.ID.String() + "/server/" + gameServer.ID.String() + "/logs",
	})

	lib.MustEncode(json.NewEncoder(w),
		gameServer)
}

func ListByHostId(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	hostId := params["id"]

	host := model.GetHost(model.DB, hostId)
	if host == nil {
		lib.MustEncode(json.NewEncoder(w),
			handler.JsonError{ErrorCode: 5024})
		return
	}

	gameServers := model.GetGameServersByHostId(model.DB, host.ID.String())
	if gameServers == nil {
		lib.MustEncode(json.NewEncoder(w),
			handler.JsonError{ErrorCode: 5025})
		return
	}

	lib.MustEncode(json.NewEncoder(w),
		gameServers)
}

func ListByUser(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user").(model.User)

	var gameServers []model.GameServer
	model.DB.Table("game_servers").Select("game_servers.*").Joins(
		"LEFT JOIN hosts ON game_servers.host_id = hosts.id").Where(
		"hosts.user_id = ?", user.ID).Find(&gameServers)

	lib.MustEncode(json.NewEncoder(w), gameServers)
}

func Get(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var gs model.GameServer

	model.DB.Where("id = ? and host_id = ?", params["server_id"], params["host_id"]).First(&gs)

	lib.MustEncode(json.NewEncoder(w), gs)
}

func ConsoleLog(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	gameServerId := params["server_id"]

	logs := model.GetLogsByGameServer(model.DB, gameServerId, 20)
	if logs == nil {
		lib.MustEncode(json.NewEncoder(w),
			handler.JsonError{ErrorCode: 5030})
		return
	}

	lib.MustEncode(json.NewEncoder(w), logs)
}

func PutLogs(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}

func Remove(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	gameServerId := params["server_id"]
	err := gs.Remove(gameServerId)
	if err != nil {
		lib.MustEncode(json.NewEncoder(w),
			handler.JsonError{ErrorCode: 5026})
		return
	}

	lib.MustEncode(json.NewEncoder(w), handler.JsonSuccess{Message: "Removing game server"})
}
