package handler

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/lib"
	"gitlab.com/systemz/aimpanel2/lib/response"
	"gitlab.com/systemz/aimpanel2/master/model"
	"gitlab.com/systemz/aimpanel2/master/service/gameserver"
	"net/http"
	"os"
	"time"
)

// swagger:route GET /host Host List
//
// List Hosts linked to the current signed-in account
//
//Responses:
//	default: jsonError
//  400: jsonError
//	200:

//TODO: find by current signed-in account
func ListHosts(w http.ResponseWriter, r *http.Request) {
	var hosts []model.Host
	user := context.Get(r, "user").(model.User)

	model.DB.Table("hosts").Where(
		"hosts.user_id = ?", user.ID).Find(&hosts)

	lib.MustEncode(json.NewEncoder(w), hosts)
}

// swagger:route GET /host/{id} Host Get
//
// Get info about Host with selected ID linked to the current signed-in account
//
//Responses:
//	default: jsonError
//  400: jsonError
//	200:

func GetHost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var host model.Host

	model.DB.Where("id = ?", params["id"]).First(&host)

	lib.MustEncode(json.NewEncoder(w), host)
}

// swagger:route POST /host Host Create
//
// Create new Host linked to the current signed-in account
//
//Responses:
//	default: jsonError
//  400: jsonError
//	200:

func CreateHost(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user").(model.User)

	host := &model.Host{}
	err := json.NewDecoder(r.Body).Decode(host)
	if err != nil {
		lib.MustEncode(json.NewEncoder(w),
			JsonError{ErrorCode: 3001, Message: "Invalid body."})
		return
	}

	host.UserId = user.ID
	host.MetricFrequency = 30

	model.DB.Save(host)

	group := model.GetGroup(model.DB, "USER-"+user.ID.String())
	if group == nil {
		lib.MustEncode(json.NewEncoder(w),
			JsonError{ErrorCode: 3002})
		return
	}

	model.DB.Save(&model.Permission{
		Name:     "Get host",
		Verb:     lib.GetVerbByName("GET"),
		GroupId:  group.ID,
		Endpoint: "/v1/host/" + host.ID.String(),
	})

	model.DB.Save(&model.Permission{
		Name:     "Delete host",
		Verb:     lib.GetVerbByName("DELETE"),
		GroupId:  group.ID,
		Endpoint: "/v1/host/" + host.ID.String(),
	})

	model.DB.Save(&model.Permission{
		Name:     "Create game server",
		Verb:     lib.GetVerbByName("POST"),
		GroupId:  group.ID,
		Endpoint: "/v1/host/" + host.ID.String() + "/server",
	})

	model.DB.Save(&model.Permission{
		Name:     "List game servers by host id",
		Verb:     lib.GetVerbByName("GET"),
		GroupId:  group.ID,
		Endpoint: "/v1/host/" + host.ID.String() + "/server",
	})

	model.DB.Save(&model.Permission{
		Name:     "Get host metric",
		Verb:     lib.GetVerbByName("GET"),
		GroupId:  group.ID,
		Endpoint: "/v1/host/" + host.ID.String() + "/metric",
	})

	lib.MustEncode(json.NewEncoder(w), host)
}

func GetHostMetric(w http.ResponseWriter, r *http.Request) {
	var metric model.MetricHost
	model.DB.Order("created_at desc").Limit(1).First(&metric)

	lib.MustEncode(json.NewEncoder(w), metric)
}

func RemoveHost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var host model.Host
	model.DB.Where("id = ?", params["id"]).First(&host)

	gameServers := model.GetGameServersByHostId(model.DB, host.ID.String())
	for _, gameServer := range *gameServers {
		err := gameserver.Remove(gameServer.ID.String())
		if err != nil {
			lib.MustEncode(json.NewEncoder(w),
				JsonError{ErrorCode: 3432})
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
			JsonError{ErrorCode: 1017})
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
			JsonError{ErrorCode: 1008})
		return
	}

	lib.MustEncode(json.NewEncoder(w), response.Token{Token: tokenString})
}

func Update(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	hostId := params["id"]

	err := gameserver.Update(hostId)
	if err != nil {
		logrus.Error(err)

		w.WriteHeader(http.StatusInternalServerError)
		lib.MustEncode(json.NewEncoder(w),
			JsonError{ErrorCode: 1234})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
