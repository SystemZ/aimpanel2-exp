package game_server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"gitlab.com/systemz/aimpanel2/lib"
	"gitlab.com/systemz/aimpanel2/lib/task"
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
	gsId := params["server_id"]

	err := gameserver.Start(gsId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		lib.MustEncode(json.NewEncoder(w),
			handler.JsonError{ErrorCode: 1234})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func Install(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	gsId := params["server_id"]

	err := gameserver.Install(gsId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		lib.MustEncode(json.NewEncoder(w),
			handler.JsonError{ErrorCode: 1234})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func Restart(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	gsId := params["server_id"]

	data := &GameServerStopReq{}
	err := json.NewDecoder(r.Body).Decode(data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		lib.MustEncode(json.NewEncoder(w),
			handler.JsonError{ErrorCode: 5012})
		return
	}

	err = gameserver.Restart(gsId, data.Type)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		lib.MustEncode(json.NewEncoder(w),
			handler.JsonError{ErrorCode: 1234})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func Stop(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	gsId := params["server_id"]

	data := &GameServerStopReq{}
	err := json.NewDecoder(r.Body).Decode(data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		lib.MustEncode(json.NewEncoder(w),
			handler.JsonError{ErrorCode: 5017})
		return
	}

	err = gameserver.Stop(gsId, data.Type)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		lib.MustEncode(json.NewEncoder(w),
			handler.JsonError{ErrorCode: 1234})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func SendCommand(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	gsId := params["server_id"]

	data := &GameServerSendCommandReq{}
	err := json.NewDecoder(r.Body).Decode(data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		lib.MustEncode(json.NewEncoder(w),
			handler.JsonError{ErrorCode: 5026})
		return
	}

	err = gameserver.SendCommand(gsId, data.Command)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		lib.MustEncode(json.NewEncoder(w),
			handler.JsonError{ErrorCode: 1234})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func Data(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	hostToken := params["host_token"]

	data := &task.Message{}
	err := json.NewDecoder(r.Body).Decode(data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
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
