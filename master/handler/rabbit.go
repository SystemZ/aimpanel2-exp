package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	rabbithole "github.com/michaelklishin/rabbit-hole"
	"github.com/sethvargo/go-password/password"
	"gitlab.com/systemz/aimpanel2/lib"
	rabbitLib "gitlab.com/systemz/aimpanel2/lib/rabbit"
	"gitlab.com/systemz/aimpanel2/master/config"
	"gitlab.com/systemz/aimpanel2/master/model"
	"gitlab.com/systemz/aimpanel2/master/rabbit"
	"net/http"
)

func GetCredentials(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var host model.Host
	model.DB.Where("token = ?", params["token"]).First(&host)

	if &host == nil {
		lib.MustEncode(json.NewEncoder(w),
			JsonError{ErrorCode: 1017})
		return
	}

	pwd, err := password.Generate(32, 10, 2, false, false)
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

	resp, err := rabbit.Client.PutUser(credentials.Username, rabbithole.UserSettings{Password: credentials.Password})
	if err != nil {
		lib.MustEncode(json.NewEncoder(w),
			JsonError{ErrorCode: 1030})
		return
	}

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		lib.MustEncode(json.NewEncoder(w),
			JsonError{ErrorCode: 1040})
		return
	}

	resp, err = rabbit.Client.UpdatePermissionsIn(credentials.VHost, credentials.Username, rabbithole.Permissions{
		Configure: ".*",
		Write:     ".*",
		Read:      ".*",
	})
	if err != nil {
		lib.MustEncode(json.NewEncoder(w),
			JsonError{ErrorCode: 1030})
		return
	}

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		lib.MustEncode(json.NewEncoder(w),
			JsonError{ErrorCode: 1041})
		return
	}

	lib.MustEncode(json.NewEncoder(w), credentials)
}
