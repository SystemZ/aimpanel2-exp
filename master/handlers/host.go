package handlers

import (
	"encoding/json"
	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"gitlab.com/systemz/aimpanel2/master/db"
	"gitlab.com/systemz/aimpanel2/master/models"
	"gitlab.com/systemz/aimpanel2/master/responses"
	"log"
	"net/http"
)

func ListHosts(w http.ResponseWriter, r *http.Request) {
	var hosts []models.Host

	db.DB.Find(&hosts)

	json.NewEncoder(w).Encode(hosts)
}

func GetHost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var host models.Host

	db.DB.First(&host, params["id"])

	json.NewEncoder(w).Encode(host)
}

func CreateHost(w http.ResponseWriter, r *http.Request) {
	host := &models.Host{}

	err := json.NewDecoder(r.Body).Decode(host)
	if err != nil {
		json.NewEncoder(w).Encode(responses.JsonError{ErrorCode: 10, Message: "Invalid body."})
		return
	}

	host.UserId = uuid.FromStringOrNil(r.Header.Get("uid"))

	log.Println(r.Header.Get("uid"))

	db.DB.Save(host)

	json.NewEncoder(w).Encode(host)
}
