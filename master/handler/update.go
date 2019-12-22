package handler

import (
	"encoding/json"
	"gitlab.com/systemz/aimpanel2/lib"
	"gitlab.com/systemz/aimpanel2/master/config"
	"gitlab.com/systemz/aimpanel2/master/model"
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
			JsonError{ErrorCode: 1234})
		return
	}

	model.Redis.Set("slave_commit", data.Commit, 0)
	model.Redis.Set("slave_url", data.Url, 0)

	w.WriteHeader(http.StatusNoContent)
}
