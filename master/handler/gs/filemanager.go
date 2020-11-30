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

// @Router /host/{hostId}/server/{gsId}/file [delete]
// @Summary Remove file
// @Tags Game Server
// @Description Remove specific file
// @Accept json
// @Produce json
// @Param hostId path string true "Host ID"
// @Param gsId path string true "Game Server ID"
// @Success 204 ""
// @Failure 400 {object} response.JsonError
// @Security ApiKey
func FileRemove(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	gsId := params["gsId"]
	oid, _ := primitive.ObjectIDFromHex(gsId)
	user := context.Get(r, "user").(model.User)

	data := &request.GameServerFile{}
	err := json.NewDecoder(r.Body).Decode(data)
	if err != nil {
		lib.ReturnError(w, http.StatusBadRequest, ecode.JsonDecode, err)
		return
	}

	err = gameserver.FileRemove(oid, data.Path, user)
	if err != nil {
		lib.ReturnError(w, http.StatusInternalServerError, ecode.GsFileRemove, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @Router /host/{hostId}/server/{gsId}/file/server [put]
// @Summary File Server
// @Tags Game Server
// @Description Start file server
// @Accept json
// @Produce json
// @Param hostId path string true "Host ID"
// @Param gsId path string true "Game Server ID"
// @Success 204 ""
// @Failure 400 {object} response.JsonError
// @Security ApiKey
func FileServer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	gsId := params["gsId"]
	oid, _ := primitive.ObjectIDFromHex(gsId)
	user := context.Get(r, "user").(model.User)

	port, err := gameserver.FileServer(oid, user)
	if err != nil {
		lib.ReturnError(w, http.StatusInternalServerError, ecode.GsFileServer, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	lib.MustEncode(json.NewEncoder(w), response.Port{Port: port})
}
