package game_server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"gitlab.com/systemz/aimpanel2/lib"
	"gitlab.com/systemz/aimpanel2/lib/task"
	"gitlab.com/systemz/aimpanel2/master/gs"
	"gitlab.com/systemz/aimpanel2/master/handler"
	"gitlab.com/systemz/aimpanel2/master/service/gameserver"
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

	if err, ok := gameserver.Start(gameServerId).(*lib.Error); ok {
		w.WriteHeader(http.StatusInternalServerError)
		lib.MustEncode(json.NewEncoder(w),
			handler.JsonError{ErrorCode: err.ErrorCode})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func Install(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	gameServerId := params["server_id"]

	if err, ok := gameserver.Install(gameServerId).(*lib.Error); ok {
		w.WriteHeader(http.StatusInternalServerError)
		lib.MustEncode(json.NewEncoder(w),
			handler.JsonError{ErrorCode: err.ErrorCode})
		return
	}

	w.WriteHeader(http.StatusNoContent)
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

func Data(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	hostToken := params["host_token"]

	data := &task.Message{}
	err := json.NewDecoder(r.Body).Decode(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		lib.MustEncode(json.NewEncoder(w),
			handler.JsonError{ErrorCode: 1234})
		return
	}

	err = gameserver.HostData(hostToken, data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		lib.MustEncode(json.NewEncoder(w),
			handler.JsonError{ErrorCode: 1234})
		return
	}

	gsId, ok := params["server_id"]
	if ok {
		err = gameserver.GsData(hostToken, gsId, data)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			lib.MustEncode(json.NewEncoder(w),
				handler.JsonError{ErrorCode: 1234})
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
