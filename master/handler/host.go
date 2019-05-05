package handler

import (
	"encoding/json"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"gitlab.com/systemz/aimpanel2/lib"
	"gitlab.com/systemz/aimpanel2/master/db"
	"gitlab.com/systemz/aimpanel2/master/model"
	"gitlab.com/systemz/aimpanel2/master/response"
	"net/http"
)

func ListHosts(w http.ResponseWriter, r *http.Request) {
	var hosts []model.Host

	db.DB.Find(&hosts)

	lib.MustEncode(json.NewEncoder(w), hosts)
}

func GetHost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var host model.Host

	db.DB.Where("id = ?", params["id"]).First(&host)

	lib.MustEncode(json.NewEncoder(w), host)
}

func CreateHost(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user").(model.User)

	host := &model.Host{}
	err := json.NewDecoder(r.Body).Decode(host)
	if err != nil {
		lib.MustEncode(json.NewEncoder(w),
			response.JsonError{ErrorCode: 3001, Message: "Invalid body."})
		return
	}

	host.UserId = user.ID
	db.DB.Save(host)

	lib.MustEncode(json.NewEncoder(w), host)
}
