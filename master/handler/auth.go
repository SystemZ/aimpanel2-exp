package handler

import (
	"encoding/json"
	"gitlab.com/systemz/aimpanel2/lib"
	"gitlab.com/systemz/aimpanel2/lib/ecode"
	"gitlab.com/systemz/aimpanel2/lib/request"
	"gitlab.com/systemz/aimpanel2/master/response"
	"gitlab.com/systemz/aimpanel2/master/service/auth"
	"net/http"
)

func Register(w http.ResponseWriter, r *http.Request) {
	data := &request.AuthRegister{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		lib.ReturnError(w, http.StatusBadRequest, ecode.JsonDecode, err)
		return
	}

	token, errCode := auth.Register(data)
	if errCode != ecode.NoError {
		lib.ReturnError(w, http.StatusBadRequest, errCode, nil)
		return
	}

	lib.MustEncode(json.NewEncoder(w), response.Token{Token: token})
}

func Login(w http.ResponseWriter, r *http.Request) {
	data := &request.AuthLogin{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		lib.ReturnError(w, http.StatusBadRequest, ecode.JsonDecode, err)
		return
	}

	token, errCode := auth.Login(data)
	if errCode != ecode.NoError {
		lib.ReturnError(w, http.StatusBadRequest, errCode, err)
		return
	}

	lib.MustEncode(json.NewEncoder(w), response.Token{Token: token})
}
