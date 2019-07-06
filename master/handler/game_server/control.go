package game_server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/lib"
	"gitlab.com/systemz/aimpanel2/master/gs"
	"gitlab.com/systemz/aimpanel2/master/handler"
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

	gameServerId := params["server_id"]

	if err, ok := gs.Start(gameServerId).(*lib.Error); ok {
		lib.MustEncode(json.NewEncoder(w),
			handler.JsonError{ErrorCode: err.ErrorCode})
		return
	}

	lib.MustEncode(json.NewEncoder(w),
		handler.JsonSuccess{Message: "Game server is starting..."})
}

func Install(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	logrus.Info("1")

	gameServerId := params["server_id"]

	if err, ok := gs.Install(gameServerId).(*lib.Error); ok {
		lib.MustEncode(json.NewEncoder(w),
			handler.JsonError{ErrorCode: err.ErrorCode})
		return
	}

	lib.MustEncode(json.NewEncoder(w),
		handler.JsonSuccess{Message: "Game server is installing."})
}

func Restart(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	gameServerId := params["server_id"]

	stopReq := &GameServerStopReq{}

	err := json.NewDecoder(r.Body).Decode(stopReq)
	if err != nil {
		lib.MustEncode(json.NewEncoder(w),
			handler.JsonError{ErrorCode: 5012})
		return
	}

	if err2, ok := gs.Restart(gameServerId, stopReq.Type).(*lib.Error); ok {
		lib.MustEncode(json.NewEncoder(w),
			handler.JsonError{ErrorCode: err2.ErrorCode})
		return
	}

	lib.MustEncode(json.NewEncoder(w), handler.JsonSuccess{Message: "Restarting the game server."})
}

func Stop(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	gameServerId := params["server_id"]

	stopReq := &GameServerStopReq{}

	err := json.NewDecoder(r.Body).Decode(stopReq)
	if err != nil {
		lib.MustEncode(json.NewEncoder(w),
			handler.JsonError{ErrorCode: 5017})
		return
	}

	if err2, ok := gs.Stop(gameServerId, stopReq.Type).(*lib.Error); ok {
		lib.MustEncode(json.NewEncoder(w),
			handler.JsonError{ErrorCode: err2.ErrorCode})
		return
	}

	lib.MustEncode(json.NewEncoder(w), handler.JsonSuccess{Message: "Stopping the game server."})
}

func SendCommand(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	gameServerId := params["server_id"]

	cmdReq := &GameServerSendCommandReq{}
	err := json.NewDecoder(r.Body).Decode(cmdReq)
	if err != nil {
		lib.MustEncode(json.NewEncoder(w),
			handler.JsonError{ErrorCode: 5026})
		return
	}

	if err2, ok := gs.SendCommand(gameServerId, cmdReq.Command).(*lib.Error); ok {
		lib.MustEncode(json.NewEncoder(w),
			handler.JsonError{ErrorCode: err2.ErrorCode})
		return
	}

	lib.MustEncode(json.NewEncoder(w), handler.JsonSuccess{Message: "Sending command to game server"})
}
