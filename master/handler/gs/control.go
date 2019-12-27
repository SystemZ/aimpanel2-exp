package gs

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"gitlab.com/systemz/aimpanel2/lib"
	"gitlab.com/systemz/aimpanel2/lib/ecode"
	"gitlab.com/systemz/aimpanel2/lib/request"
	"gitlab.com/systemz/aimpanel2/lib/task"
	"gitlab.com/systemz/aimpanel2/master/handler"
	"gitlab.com/systemz/aimpanel2/master/service/gameserver"
	"log"
	"net/http"
)

// @Router /host/{host_id}/server/{server_id}/start [put]
// @Summary Start
// @Tags Game Server
// @Description Start selected game server
// @Accept json
// @Produce json
// @Param host_id path string true "Host ID"
// @Param server_id path string true "Game Server ID"
// @Success 204 ""
// @Failure 400 {object} handler.JsonError
// @Security ApiKey
func Start(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	gsId := params["server_id"]

	err := gameserver.Start(gsId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		lib.MustEncode(json.NewEncoder(w),
			handler.JsonError{ErrorCode: ecode.GsStart})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @Router /host/{host_id}/server/{server_id}/install [put]
// @Summary Install
// @Tags Game Server
// @Description Install selected game server
// @Accept json
// @Produce json
// @Param host_id path string true "Host ID"
// @Param server_id path string true "Game Server ID"
// @Success 204 ""
// @Failure 400 {object} handler.JsonError
// @Security ApiKey
func Install(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	gsId := params["server_id"]

	err := gameserver.Install(gsId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		lib.MustEncode(json.NewEncoder(w),
			handler.JsonError{ErrorCode: ecode.GsInstall})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @Router /host/{host_id}/server/{server_id}/restart [put]
// @Summary Restart
// @Tags Game Server
// @Description Restart selected game server
// @Accept json
// @Produce json
// @Param host_id path string true "Host ID"
// @Param server_id path string true "Game Server ID"
// @Success 204 ""
// @Failure 400 {object} handler.JsonError
// @Security ApiKey
func Restart(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	gsId := params["server_id"]

	data := &request.GameServerStop{}
	err := json.NewDecoder(r.Body).Decode(data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		lib.MustEncode(json.NewEncoder(w),
			handler.JsonError{ErrorCode: ecode.JsonDecode})
		return
	}

	err = gameserver.Restart(gsId, data.Type)
	if err != nil {
		log.Printf("fail restarting GS: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		lib.MustEncode(json.NewEncoder(w),
			handler.JsonError{ErrorCode: ecode.GsRestart})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @Router /host/{host_id}/server/{server_id}/stop [put]
// @Summary Stop
// @Tags Game Server
// @Description Stop selected game server
// @Accept json
// @Produce json
// @Param host_id path string true "Host ID"
// @Param server_id path string true "Game Server ID"
// @Success 204 ""
// @Failure 400 {object} handler.JsonError
// @Security ApiKey
func Stop(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	gsId := params["server_id"]

	data := &request.GameServerStop{}
	err := json.NewDecoder(r.Body).Decode(data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		lib.MustEncode(json.NewEncoder(w),
			handler.JsonError{ErrorCode: ecode.JsonDecode})
		return
	}

	err = gameserver.Stop(gsId, data.Type)
	if err != nil {
		log.Printf("fail stopping GS: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		lib.MustEncode(json.NewEncoder(w),
			handler.JsonError{ErrorCode: ecode.GsStop})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @Router /host/{host_id}/server/{server_id}/command [put]
// @Summary Send command
// @Tags Game Server
// @Description Send command to selected game server
// @Accept json
// @Produce json
// @Param host_id path string true "Host ID"
// @Param server_id path string true "Game Server ID"
// @Param host body request.GameServerSendCommand true " "
// @Success 204 ""
// @Failure 400 {object} handler.JsonError
// @Security ApiKey
func SendCommand(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	gsId := params["server_id"]

	data := &request.GameServerSendCommand{}
	err := json.NewDecoder(r.Body).Decode(data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		lib.MustEncode(json.NewEncoder(w),
			handler.JsonError{ErrorCode: ecode.JsonDecode})
		return
	}

	err = gameserver.SendCommand(gsId, data.Command)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		lib.MustEncode(json.NewEncoder(w),
			handler.JsonError{ErrorCode: ecode.GsCmd})
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
			handler.JsonError{ErrorCode: ecode.JsonDecode})
		return
	}

	err = gameserver.HostData(hostToken, data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		lib.MustEncode(json.NewEncoder(w),
			handler.JsonError{ErrorCode: ecode.HostData})
		return
	}

	gsId, ok := params["server_id"]
	if ok {
		err = gameserver.GsData(hostToken, gsId, data)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			lib.MustEncode(json.NewEncoder(w),
				handler.JsonError{ErrorCode: ecode.GsData})
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
