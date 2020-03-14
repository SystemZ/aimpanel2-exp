package handler

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/lib"
	"gitlab.com/systemz/aimpanel2/lib/ecode"
	"gitlab.com/systemz/aimpanel2/master/config"
	"gitlab.com/systemz/aimpanel2/master/model"
	"gitlab.com/systemz/aimpanel2/master/service/gameserver"
	"net/http"
)

type NewVersionReq struct {
	Commit string `json:"commit"`
	Url    string `json:"url"`
}

func NewVersion(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Token") != config.UPDATE_TOKEN {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	data := &NewVersionReq{}
	err := json.NewDecoder(r.Body).Decode(data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		lib.MustEncode(json.NewEncoder(w),
			JsonError{ErrorCode: ecode.JsonDecode})
		return
	}

	err = model.Redis.Set("slave_commit", data.Commit, 0).Err()
	if err != nil {
		logrus.Error(err)
	}

	err = model.Redis.Set("slave_url", data.Url, 0).Err()
	if err != nil {
		logrus.Error(err)
	}

	//TODO: Find better way to update all hosts
	hosts := model.GetHosts()
	go func() {
		for _, host := range hosts {
			err := gameserver.Update(host.ID)
			if err != nil {
				logrus.Error(err)
			}
		}
	}()

	w.WriteHeader(http.StatusNoContent)
}
