package game_server

import (
	"encoding/json"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"gitlab.com/systemz/aimpanel2/lib"
	"gitlab.com/systemz/aimpanel2/master/gs"
	"gitlab.com/systemz/aimpanel2/master/handler"
	"gitlab.com/systemz/aimpanel2/master/model"

	"net/http"
)

type GameServerStopReq struct {
	Type uint `json:"type"`
}

type GameServerSendCommandReq struct {
	Command string `json:"command"`
}

func Start(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	hostId := params["host_id"]
	gameServerId := params["server_id"]

	user := context.Get(r, "user").(model.User)

	err := gs.Start(hostId, gameServerId, user).(*lib.Error)
	if err != nil {
		lib.MustEncode(json.NewEncoder(w),
			handler.JsonError{ErrorCode: err.ErrorCode})
	}

	lib.MustEncode(json.NewEncoder(w),
		handler.JsonSuccess{Message: "Game server is starting..."})
}

func Install(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	hostId := params["host_id"]
	gameServerId := params["server_id"]
	user := context.Get(r, "user").(model.User)

	err := gs.Install(hostId, gameServerId, user).(*lib.Error)
	if err != nil {

	}

	lib.MustEncode(json.NewEncoder(w),
		handler.JsonSuccess{Message: "Game server is installing."})
}

func Restart(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	hostId := params["host_id"]
	gameServerId := params["server_id"]
	user := context.Get(r, "user").(model.User)

	stopReq := &GameServerStopReq{}

	err := json.NewDecoder(r.Body).Decode(stopReq)
	if err != nil {
		lib.MustEncode(json.NewEncoder(w),
			handler.JsonError{ErrorCode: 5012})
		return
	}

	err2 := gs.Restart(hostId, gameServerId, user, stopReq.Type).(*lib.Error)
	if err2 != nil {
		lib.MustEncode(json.NewEncoder(w),
			handler.JsonError{ErrorCode: err2.ErrorCode})
	}

	lib.MustEncode(json.NewEncoder(w), handler.JsonSuccess{Message: "Restarting the game server."})
}

func Stop(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	hostId := params["host_id"]
	gameServerId := params["server_id"]
	user := context.Get(r, "user").(model.User)

	stopReq := &GameServerStopReq{}

	err := json.NewDecoder(r.Body).Decode(stopReq)
	if err != nil {
		lib.MustEncode(json.NewEncoder(w),
			handler.JsonError{ErrorCode: 5017})
		return
	}

	err2 := gs.Stop(hostId, gameServerId, user, stopReq.Type).(*lib.Error)
	if err2 != nil {
		lib.MustEncode(json.NewEncoder(w),
			handler.JsonError{ErrorCode: err2.ErrorCode})
	}

	lib.MustEncode(json.NewEncoder(w), handler.JsonSuccess{Message: "Stopping the game server."})
}

func SendCommand(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	hostId := params["host_id"]
	gameServerId := params["server_id"]
	user := context.Get(r, "user").(model.User)

	cmdReq := &GameServerSendCommandReq{}
	err := json.NewDecoder(r.Body).Decode(cmdReq)
	if err != nil {
		lib.MustEncode(json.NewEncoder(w),
			handler.JsonError{ErrorCode: 5026})
		return
	}

	err2 := gs.SendCommand(hostId, gameServerId, user, cmdReq.Command).(*lib.Error)
	if err2 != nil {
		lib.MustEncode(json.NewEncoder(w),
			handler.JsonError{ErrorCode: err2.ErrorCode})
	}

	lib.MustEncode(json.NewEncoder(w), handler.JsonSuccess{Message: "Sending command to game server"})
}
