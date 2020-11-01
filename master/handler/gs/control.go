package gs

import (
	"encoding/json"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"gitlab.com/systemz/aimpanel2/lib"
	"gitlab.com/systemz/aimpanel2/lib/ecode"
	"gitlab.com/systemz/aimpanel2/lib/request"
	"gitlab.com/systemz/aimpanel2/master/model"
	"gitlab.com/systemz/aimpanel2/master/response"
	"gitlab.com/systemz/aimpanel2/master/service/gameserver"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
// @Failure 400 {object} response.JsonError
// @Security ApiKey
func Start(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	gsId := params["gsId"]
	oid, _ := primitive.ObjectIDFromHex(gsId)
	user := context.Get(r, "user").(model.User)

	err := gameserver.Start(oid, user)
	if err != nil {
		lib.ReturnError(w, http.StatusInternalServerError, ecode.GsStart, err)
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
// @Failure 400 {object} response.JsonError
// @Security ApiKey
func Install(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	gsId := params["gsId"]
	oid, err := primitive.ObjectIDFromHex(gsId)
	user := context.Get(r, "user").(model.User)
	if err != nil {
		lib.ReturnError(w, http.StatusBadRequest, ecode.OidError, err)
		return
	}

	err = gameserver.Install(oid, user)
	if err != nil {
		lib.ReturnError(w, http.StatusInternalServerError, ecode.GsInstall, err)
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
// @Failure 400 {object} response.JsonError
// @Security ApiKey
func Restart(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	gsId := params["gsId"]
	oid, _ := primitive.ObjectIDFromHex(gsId)
	user := context.Get(r, "user").(model.User)

	data := &request.GameServerStop{}
	err := json.NewDecoder(r.Body).Decode(data)
	if err != nil {
		lib.ReturnError(w, http.StatusBadRequest, ecode.JsonDecode, err)
		return
	}

	err = gameserver.Restart(oid, data.Type, user)
	if err != nil {
		lib.ReturnError(w, http.StatusInternalServerError, ecode.GsRestart, err)
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
// @Failure 400 {object} response.JsonError
// @Security ApiKey
func Stop(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	gsId := params["gsId"]
	oid, _ := primitive.ObjectIDFromHex(gsId)
	user := context.Get(r, "user").(model.User)

	data := &request.GameServerStop{}
	err := json.NewDecoder(r.Body).Decode(data)
	if err != nil {
		lib.ReturnError(w, http.StatusBadRequest, ecode.JsonDecode, err)
		return
	}

	err = gameserver.Stop(oid, data.Type, user)
	if err != nil {
		lib.ReturnError(w, http.StatusInternalServerError, ecode.GsStop, err)
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
// @Failure 400 {object} response.JsonError
// @Security ApiKey
func SendCommand(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	gsId := params["gsId"]
	oid, _ := primitive.ObjectIDFromHex(gsId)
	user := context.Get(r, "user").(model.User)

	data := &request.GameServerSendCommand{}
	err := json.NewDecoder(r.Body).Decode(data)
	if err != nil {
		lib.ReturnError(w, http.StatusBadRequest, ecode.JsonDecode, err)
		return
	}

	err = gameserver.SendCommand(oid, data.Command, user)
	if err != nil {
		lib.ReturnError(w, http.StatusInternalServerError, ecode.GsCmd, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @Router /host/{hostId}/server/{gsId}/file/list [get]
// @Summary File List
// @Tags Game Server
// @Description Get game server file list
// @Accept json
// @Produce json
// @Param hostId path string true "Host ID"
// @Param gsId path string true "Game Server ID"
// @Success 200 {object} filemanager.Node
// @Failure 400 {object} response.JsonError
// @Security ApiKey
func FileList(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	gsId := params["gsId"]
	oid, _ := primitive.ObjectIDFromHex(gsId)
	user := context.Get(r, "user").(model.User)

	files, err := gameserver.FileList(oid, user)
	if err != nil {
		lib.ReturnError(w, http.StatusInternalServerError, ecode.GsFileList, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	lib.MustEncode(json.NewEncoder(w), files)
}

// @Router /host/{host_id}/server/{server_id}/shutdown [put]
// @Summary Shutdown
// @Tags Game Server
// @Description Shutdown selected game server
// @Accept json
// @Produce json
// @Param host_id path string true "Host ID"
// @Param server_id path string true "Game Server ID"
// @Success 204 ""
// @Failure 400 {object} response.JsonError
// @Security ApiKey
func Shutdown(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	gsId := params["gsId"]
	oid, err := primitive.ObjectIDFromHex(gsId)
	user := context.Get(r, "user").(model.User)
	if err != nil {
		lib.ReturnError(w, http.StatusBadRequest, ecode.OidError, err)
		return
	}

	err = gameserver.Shutdown(oid, user)
	if err != nil {
		lib.ReturnError(w, http.StatusInternalServerError, ecode.GsInstall, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @Router /host/{host_id}/server/{server_id}/backup [put]
// @Summary Backup
// @Tags Game Server
// @Description Backup selected game server
// @Accept json
// @Produce json
// @Param host_id path string true "Host ID"
// @Param server_id path string true "Game Server ID"
// @Success 204 ""
// @Failure 400 {object} response.JsonError
// @Security ApiKey
func Backup(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	gsId := params["gsId"]
	oid, err := primitive.ObjectIDFromHex(gsId)
	user := context.Get(r, "user").(model.User)
	if err != nil {
		lib.ReturnError(w, http.StatusBadRequest, ecode.OidError, err)
		return
	}

	err = gameserver.Backup(oid, user)
	if err != nil {
		lib.ReturnError(w, http.StatusInternalServerError, ecode.GsBackup, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @Router /host/{hostId}/server/{gsId}/backup/list [get]
// @Summary Backup List
// @Tags Game Server
// @Description Get game server backup list
// @Accept json
// @Produce json
// @Param hostId path string true "Host ID"
// @Param gsId path string true "Game Server ID"
// @Success 200 {object} response.BackupList
// @Failure 400 {object} response.JsonError
// @Security ApiKey
func BackupList(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	gsId := params["gsId"]
	oid, _ := primitive.ObjectIDFromHex(gsId)
	user := context.Get(r, "user").(model.User)

	files, err := gameserver.BackupList(oid, user)
	if err != nil {
		lib.ReturnError(w, http.StatusInternalServerError, ecode.GsBackupList, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	lib.MustEncode(json.NewEncoder(w), response.BackupList{Backups: files})
}

// @Router /host/{host_id}/server/{server_id}/backup/restore [put]
// @Summary Backup restore
// @Tags Game Server
// @Description Restore specific backup for gs
// @Accept json
// @Produce json
// @Param host_id path string true "Host ID"
// @Param server_id path string true "Game Server ID"
// @Param host body request.BackupRestore true " "
// @Success 204 ""
// @Failure 400 {object} response.JsonError
// @Security ApiKey
func BackupRestore(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	gsId := params["gsId"]
	oid, _ := primitive.ObjectIDFromHex(gsId)
	user := context.Get(r, "user").(model.User)

	data := &request.BackupRestore{}
	err := json.NewDecoder(r.Body).Decode(data)
	if err != nil {
		lib.ReturnError(w, http.StatusBadRequest, ecode.JsonDecode, err)
		return
	}

	err = gameserver.BackupRestore(oid, data.BackupFilename, user)
	if err != nil {
		lib.ReturnError(w, http.StatusInternalServerError, ecode.GsBackupRestore, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
