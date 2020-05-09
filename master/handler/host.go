package handler

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
	"gitlab.com/systemz/aimpanel2/master/service/host"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

// @Router /host [get]
// @Summary List
// @Tags Host
// @Description List Hosts linked to the current signed-in account
// @Accept json
// @Produce json
// @Success 200 {object} response.HostList
// @Failure 400 {object} response.JsonError
// @Security ApiKey
func HostList(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user").(model.User)
	hosts, err := model.GetHostsByUserId(user.ID)
	if err != nil {
		lib.ReturnError(w, http.StatusInternalServerError, ecode.DbError, nil)
		return
	}

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
// @Failure 400 {object} response.JsonError
// @Security ApiKey
func HostDetails(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	oid, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		lib.ReturnError(w, http.StatusBadRequest, ecode.OidError, nil)
		return
	}

	h, err := model.GetHostById(oid)
	if err != nil {
		lib.ReturnError(w, http.StatusInternalServerError, ecode.DbError, nil)
		return
	}

	if h == nil {
		lib.ReturnError(w, http.StatusBadRequest, ecode.HostNotFound, nil)
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
// @Failure 400 {object} response.JsonError
// @Security ApiKey
func HostCreate(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user").(model.User)

	data := &request.HostCreate{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		lib.ReturnError(w, http.StatusBadRequest, ecode.JsonDecode, err)
		return
	}

	h, errCode := host.Create(data, user.ID)
	if errCode != ecode.NoError {
		lib.ReturnError(w, http.StatusInternalServerError, errCode, nil)
		return
	}

	lib.MustEncode(json.NewEncoder(w), response.Token{Token: h.Token})
}

// @Router /host/{id}/metric [get]
// @Summary Metric
// @Tags Host
// @Description Get latest metrics for host with ID, linked to the current signed-in account
// @Accept json
// @Produce json
// @Param id path string true "Host ID"
// @Success 200 {object} response.HostMetrics
// @Failure 400 {object} response.JsonError
// @Security ApiKey
func HostMetric(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	oid, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		lib.ReturnError(w, http.StatusBadRequest, ecode.OidError, nil)
		return
	}

	metrics, err := model.GetHostMetricsByHostId(oid, 20)
	if err != nil {
		lib.ReturnError(w, http.StatusInternalServerError, ecode.DbError, nil)
		return
	}

	lib.MustEncode(json.NewEncoder(w), response.HostMetrics{Metrics: metrics})
}

// @Router /host/{id} [delete]
// @Summary Remove
// @Tags Host
// @Description Removes host with all linked game servers
// @Accept json
// @Produce json
// @Param id path string true "Host ID"
// @Success 200 {object} response.JsonSuccess
// @Failure 400 {object} response.JsonError
// @Security ApiKey
func HostRemove(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	oid, _ := primitive.ObjectIDFromHex(params["id"])
	errCode := host.Remove(oid)
	if errCode != ecode.NoError {
		lib.ReturnError(w, http.StatusBadRequest, errCode, nil)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func HostAuth(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	token, errCode := host.Auth(params["token"])
	if errCode != ecode.NoError {
		lib.ReturnError(w, http.StatusBadRequest, errCode, nil)
		return
	}

	lib.MustEncode(json.NewEncoder(w), response.Token{Token: token})
}

//TODO: Available for users?
func HostUpdate(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	oid, _ := primitive.ObjectIDFromHex(params["id"])
	err := gameserver.Update(oid)
	if err != nil {
		lib.ReturnError(w, http.StatusInternalServerError, ecode.GsUpdate, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @Router /host/{id}/job [post]
// @Summary Create Job
// @Tags Host
// @Description Create new Job linked to the given host
// @Accept json
// @Produce json
// @Param host body request.HostCreateJob true " "
// @Success 200 {object} response.Token
// @Failure 400 {object} response.JsonError
// @Security ApiKey
func HostCreateJob(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user").(model.User)

	params := mux.Vars(r)

	data := &request.HostCreateJob{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		lib.ReturnError(w, http.StatusBadRequest, ecode.JsonDecode, err)
		return
	}

	oid, _ := primitive.ObjectIDFromHex(params["id"])
	_, errCode := host.CreateJob(data, user.ID, oid)
	if errCode != ecode.NoError {
		lib.ReturnError(w, http.StatusInternalServerError, errCode, nil)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @Router /host/{id}/job/{job_id} [delete]
// @Summary Remove Job
// @Tags Host
// @Description Removes host job
// @Accept json
// @Produce json
// @Param id path string true "Host ID"
// @Param job_id path string true "Job ID"
// @Success 200 {object} response.JsonSuccess
// @Failure 400 {object} response.JsonError
// @Security ApiKey
func HostJobRemove(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	hostId, _ := primitive.ObjectIDFromHex(params["id"])
	jobId, _ := primitive.ObjectIDFromHex(params["job_id"])
	errCode := host.RemoveJob(hostId, jobId)
	if errCode != ecode.NoError {
		lib.ReturnError(w, http.StatusBadRequest, errCode, nil)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func HostJobList(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	oid, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		lib.ReturnError(w, http.StatusBadRequest, ecode.OidError, nil)
		return
	}

	jobs, err := model.GetHostJobsByHostId(oid)
	if err != nil {
		lib.ReturnError(w, http.StatusInternalServerError, ecode.DbError, nil)
		return
	}

	lib.MustEncode(json.NewEncoder(w), response.HostJobList{Jobs: jobs})
}
