package handler

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/lib"
	"gitlab.com/systemz/aimpanel2/lib/ecode"
	"gitlab.com/systemz/aimpanel2/lib/request"
	"gitlab.com/systemz/aimpanel2/master/model"
	"gitlab.com/systemz/aimpanel2/master/response"
	"gitlab.com/systemz/aimpanel2/master/service/gameserver"
	"net/http"
	"os"
	"time"
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
	var hosts []model.Host
	user := context.Get(r, "user").(model.User)

	model.DB.Table("hosts").Where(
		"hosts.user_id = ?", user.ID).Find(&hosts)

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

	var host model.Host

	model.DB.Where("id = ?", params["id"]).First(&host)

	lib.MustEncode(json.NewEncoder(w), response.Host{Host: host})
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

	host := &model.Host{
		Name:            data.Name,
		Ip:              data.Ip,
		UserId:          user.ID,
		MetricFrequency: 30,
	}
	model.DB.Save(&host)

	group := model.GetGroup(model.DB, "USER-"+user.ID.String())
	if group == nil {
		lib.MustEncode(json.NewEncoder(w),
			JsonError{ErrorCode: ecode.GroupNotFound})
		return
	}

	// FIXME handle errors
	model.CreatePermissionsForNewHost(group.ID, host.ID.String())

	lib.MustEncode(json.NewEncoder(w), response.Token{Token: host.Token})
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

	var metrics []model.MetricHost
	model.DB.Where("host_id = ?", params["id"]).Order("created_at desc").Limit(1).Find(&metrics)

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

	var host model.Host
	model.DB.Where("id = ?", params["id"]).First(&host)

	gameServers := model.GetGameServersByHostId(model.DB, host.ID.String())
	for _, gameServer := range *gameServers {
		err := gameserver.Remove(gameServer.ID.String())
		if err != nil {
			lib.MustEncode(json.NewEncoder(w),
				JsonError{ErrorCode: ecode.GsRemove})
			return
		}
	}

	model.DB.Where("endpoint LIKE ?", "/v1/host/"+host.ID.String()+"%").Delete(&model.Permission{})
	model.DB.Delete(&host)

	lib.MustEncode(json.NewEncoder(w), JsonSuccess{Message: "Removing host"})
}

func HostAuth(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var host model.Host
	model.DB.Where("token = ?", params["token"]).First(&host)

	if &host == nil {
		lib.MustEncode(json.NewEncoder(w),
			JsonError{ErrorCode: ecode.HostNotFound})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * 48).Unix(),
		"uid": host.ID.String(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		lib.MustEncode(json.NewEncoder(w),
			JsonError{ErrorCode: ecode.Unknown})
		return
	}

	lib.MustEncode(json.NewEncoder(w), response.Token{Token: tokenString})
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
