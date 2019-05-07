package game_server

import (
	"encoding/json"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"gitlab.com/systemz/aimpanel2/lib"
	"gitlab.com/systemz/aimpanel2/master/db"
	"gitlab.com/systemz/aimpanel2/master/model"
	"gitlab.com/systemz/aimpanel2/master/response"
	"net/http"
)

func Create(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	hostId := params["host_id"]

	gameServer := &model.GameServer{}
	err := json.NewDecoder(r.Body).Decode(gameServer)
	if err != nil {
		lib.MustEncode(json.NewEncoder(w),
			response.JsonError{ErrorCode: 5022})
		return
	}

	user := context.Get(r, "user").(model.User)

	host := user.GetHost(db.DB, hostId)
	if host == nil {
		lib.MustEncode(json.NewEncoder(w),
			response.JsonError{ErrorCode: 5023})
		return
	}

	gameServer.HostId = host.ID

	db.DB.Save(gameServer)

	lib.MustEncode(json.NewEncoder(w),
		gameServer)
}

func ListByHostId(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	hostId := params["id"]

	user := context.Get(r, "user").(model.User)

	host := user.GetHost(db.DB, hostId)
	if host == nil {
		lib.MustEncode(json.NewEncoder(w),
			response.JsonError{ErrorCode: 5024})
		return
	}

	gameServers := host.GetGameServers(db.DB)
	if gameServers == nil {
		lib.MustEncode(json.NewEncoder(w),
			response.JsonError{ErrorCode: 5025})
		return
	}

	lib.MustEncode(json.NewEncoder(w),
		gameServers)
}

func ListByUser(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user").(model.User)

	var gameServers []model.GameServer
	db.DB.Table("game_servers").Select("game_servers.*").Joins(
		"LEFT JOIN hosts ON game_servers.host_id = hosts.id").Where(
		"hosts.user_id = ?", user.ID).Find(&gameServers)

	lib.MustEncode(json.NewEncoder(w), gameServers)
}
