package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"gitlab.com/systemz/aimpanel2/lib"
	"gitlab.com/systemz/aimpanel2/lib/ecode"
	"gitlab.com/systemz/aimpanel2/lib/task"
	"gitlab.com/systemz/aimpanel2/master/service/gameserver"
	"gitlab.com/systemz/aimpanel2/master/service/host"
	"net/http"
)

func ReceiveData(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	hostToken := params["hostToken"]

	data := &task.Message{}
	err := json.NewDecoder(r.Body).Decode(data)
	if err != nil {
		lib.ReturnError(w, http.StatusBadRequest, ecode.JsonDecode, err)
		return
	}

	if data.GameServerID != "" {
		err = gameserver.Data(hostToken, data)
		if err != nil {
			lib.ReturnError(w, http.StatusInternalServerError, ecode.GsData, err)
			return
		}
	} else {
		err = host.Data(hostToken, data)
		if err != nil {
			lib.ReturnError(w, http.StatusInternalServerError, ecode.HostData, err)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

func ReceiveBatchData(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	hostToken := params["hostToken"]

	var taskMsgs task.Messages
	err := json.NewDecoder(r.Body).Decode(&taskMsgs)
	if err != nil {
		lib.ReturnError(w, http.StatusBadRequest, ecode.JsonDecode, err)
		return
	}

	for _, taskMsg := range taskMsgs {
		if taskMsg.GameServerID != "" {
			err = gameserver.Data(hostToken, &taskMsg)
			if err != nil {
				lib.ReturnError(w, http.StatusInternalServerError, ecode.GsData, err)
				return
			}
		} else {
			err = host.Data(hostToken, &taskMsg)
			if err != nil {
				lib.ReturnError(w, http.StatusInternalServerError, ecode.HostData, err)
				return
			}
		}
	}
	w.WriteHeader(http.StatusNoContent)
}
