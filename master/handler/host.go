package handler

import (
	"encoding/json"
	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"gitlab.com/systemz/aimpanel2/master/db"
	"gitlab.com/systemz/aimpanel2/master/model"
	"gitlab.com/systemz/aimpanel2/master/response"
	"log"
	"net/http"
)

func ListHosts(w http.ResponseWriter, r *http.Request) {
	// swagger:route POST /hosts Host list
	//
	// List hosts
	//
	//Consumes:
	//	- application/json
	//
	//Produces:
	//	- application/json
	//
	//Schemes: http, https
	//
	//Responses:
	//	default: jsonError
	//	200:
	var hosts []model.Host

	db.DB.Find(&hosts)

	json.NewEncoder(w).Encode(hosts)
}

func GetHost(w http.ResponseWriter, r *http.Request) {
	// swagger:route POST /hosts/{id} Host get
	//
	// Get a host
	//
	//Consumes:
	//	- application/json
	//
	//Produces:
	//	- application/json
	//
	//Schemes: http, https
	//
	//Responses:
	//	default: jsonError
	//	200:
	params := mux.Vars(r)

	var host model.Host

	db.DB.Where("id = ?", params["id"]).First(&host)

	json.NewEncoder(w).Encode(host)
}

func CreateHost(w http.ResponseWriter, r *http.Request) {
	// swagger:route POST /hosts Host create
	//
	// Creates new host
	//
	//Consumes:
	//	- application/json
	//
	//Produces:
	//	- application/json
	//
	//Schemes: http, https
	//
	//Responses:
	//	default: jsonError
	//	200:
	host := &model.Host{}

	err := json.NewDecoder(r.Body).Decode(host)
	if err != nil {
		json.NewEncoder(w).Encode(response.JsonError{ErrorCode: 10, Message: "Invalid body."})
		return
	}

	host.UserId = uuid.FromStringOrNil(r.Header.Get("uid"))

	log.Println(r.Header.Get("uid"))

	db.DB.Save(host)

	json.NewEncoder(w).Encode(host)
}
