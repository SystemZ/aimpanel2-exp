package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sethvargo/go-password/password"
	"gitlab.com/systemz/aimpanel2/lib"
	rabbitLib "gitlab.com/systemz/aimpanel2/lib/rabbit"
	"gitlab.com/systemz/aimpanel2/master/config"
	"gitlab.com/systemz/aimpanel2/master/model"
	"gitlab.com/systemz/aimpanel2/master/rabbit"
	"net/http"
)

func GetHostCredentials(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var host model.Host
	model.DB.Where("token = ?", params["token"]).First(&host)

	if &host == nil {
		lib.MustEncode(json.NewEncoder(w),
			JsonError{ErrorCode: 1017})
		return
	}

	pwd, err := password.Generate(32, 10, 0, false, false)
	if err != nil {
		lib.MustEncode(json.NewEncoder(w),
			JsonError{ErrorCode: 1019})
		return
	}

	credentials := rabbitLib.Credentials{
		Host:     config.RABBITMQ_HOST,
		Port:     config.RABBITMQ_PORT,
		Username: host.ID.String(),
		Password: pwd,
		VHost:    config.RABBITMQ_VHOST,
	}

	err = rabbit.PutUser(credentials)
	if err != nil {
		lib.MustEncode(json.NewEncoder(w),
			JsonError{ErrorCode: 1020})
		return
	}

	lib.MustEncode(json.NewEncoder(w), credentials)
}

func GetGameServerCredentials(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var host model.Host
	model.DB.Where("token = ?", params["token"]).First(&host)

	if &host == nil {
		lib.MustEncode(json.NewEncoder(w),
			JsonError{ErrorCode: 1017})
		return
	}

	var gameServer model.GameServer
	model.DB.Where("id = ? ", params["server_id"]).First(&gameServer)
	if &gameServer == nil {
		lib.MustEncode(json.NewEncoder(w),
			JsonError{ErrorCode: 1018})
		return
	}

	pwd, err := password.Generate(32, 10, 0, false, false)
	if err != nil {
		lib.MustEncode(json.NewEncoder(w),
			JsonError{ErrorCode: 1019})
		return
	}

	credentials := rabbitLib.Credentials{
		Host:     config.RABBITMQ_HOST,
		Port:     config.RABBITMQ_PORT,
		Username: gameServer.ID.String(),
		Password: pwd,
		VHost:    config.RABBITMQ_VHOST,
	}

	err = rabbit.PutUser(credentials)
	if err != nil {
		lib.MustEncode(json.NewEncoder(w),
			JsonError{ErrorCode: 1020})
		return
	}

	lib.MustEncode(json.NewEncoder(w), credentials)
}
