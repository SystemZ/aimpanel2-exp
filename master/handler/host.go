package handler

import (
	"encoding/json"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/lib"
	"gitlab.com/systemz/aimpanel2/lib/ecode"
	"gitlab.com/systemz/aimpanel2/lib/request"
	"gitlab.com/systemz/aimpanel2/master/model"
	"gitlab.com/systemz/aimpanel2/master/response"
	"gitlab.com/systemz/aimpanel2/master/service/gameserver"
	"gitlab.com/systemz/aimpanel2/master/service/host"
	"net/http"
)

// @Router /host [get]
// @Summary List
// @Tags Host
// @Description List Hosts linked to the current signed-in account
// @Accept json
// @Produce json
// @Success 200 {object} response.HostList
// @Failure 400 {object} JsonError
// @Security ApiKey
func HostList(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user").(model.User)
	hosts := model.GetHostsByUserId(model.DB, user.ID)
	lib.MustEncode(json.NewEncoder(w), response.HostList{Hosts: hosts})
}

// @Router /host/{id} [get]
// @Summary Details
// @Tags Host
// @Description Get details about Host with selected ID linked to the current signed-in account
// @Accept json
// @Produce json
// @Param id path string true "Host ID"
// @Success 200 {object} response.Host
// @Failure 400 {object} JsonError
// @Security ApiKey
func HostDetails(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	h := model.GetHost(model.DB, params["id"])
	if h == nil {
		w.WriteHeader(http.StatusBadRequest)
		lib.MustEncode(json.NewEncoder(w),
			JsonError{ErrorCode: ecode.HostNotFound})
		return
	}

	lib.MustEncode(json.NewEncoder(w), response.Host{Host: *h})
}

// @Router /host [post]
// @Summary Create
// @Tags Host
// @Description Create new Host linked to the current signed-in account
// @Accept json
// @Produce json
// @Param host body request.HostCreate true " "
// @Success 200 {object} response.Token
// @Failure 400 {object} JsonError
// @Security ApiKey
func HostCreate(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user").(model.User)

	data := &request.HostCreate{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		lib.MustEncode(json.NewEncoder(w),
			JsonError{ErrorCode: ecode.JsonDecode})
		return
	}

	h, errCode := host.Create(data, user.ID)
	if errCode != ecode.NoError {
		lib.MustEncode(json.NewEncoder(w),
			JsonError{ErrorCode: errCode})
		return
	}

	lib.MustEncode(json.NewEncoder(w), response.Token{Token: h.Token})
}

// @Router /host/{id}/metric [get]
// @Summary Metric
// @Tags Host
// @Description Get last host metric with selected ID linked to the current signed-in account
// @Accept json
// @Produce json
// @Param id path string true "Host ID"
// @Success 200 {object} response.HostMetrics
// @Failure 400 {object} JsonError
// @Security ApiKey
func HostMetric(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	metrics := model.GetHostMetrics(model.DB, params["id"], 1)
	lib.MustEncode(json.NewEncoder(w), response.HostMetrics{Metrics: metrics})
}

// @Router /host/{id} [delete]
// @Summary Remove
// @Tags Host
// @Description Removes host with all linked game servers
// @Accept json
// @Produce json
// @Param id path string true "Host ID"
// @Success 200 {object} JsonSuccess
// @Failure 400 {object} JsonError
// @Security ApiKey
func HostRemove(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	errCode := host.Remove(params["id"])
	if errCode != ecode.NoError {
		w.WriteHeader(http.StatusBadRequest)
		lib.MustEncode(json.NewEncoder(w),
			JsonError{ErrorCode: errCode})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func HostAuth(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	token, errCode := host.Auth(params["token"])
	if errCode != ecode.NoError {
		w.WriteHeader(http.StatusBadRequest)
		lib.MustEncode(json.NewEncoder(w),
			JsonError{ErrorCode: errCode})
		return
	}

	lib.MustEncode(json.NewEncoder(w), response.Token{Token: token})
}

//TODO: Available for users?
func HostUpdate(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	hostId := params["id"]

	err := gameserver.Update(hostId)
	if err != nil {
		logrus.Error(err)

		w.WriteHeader(http.StatusInternalServerError)
		lib.MustEncode(json.NewEncoder(w),
			JsonError{ErrorCode: ecode.GsUpdate})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
