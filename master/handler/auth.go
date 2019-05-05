package handler

import (
	"encoding/json"
	"gitlab.com/systemz/aimpanel2/lib"
	"gitlab.com/systemz/aimpanel2/master/db"
	"gitlab.com/systemz/aimpanel2/master/model"
	"gitlab.com/systemz/aimpanel2/master/request"
	"gitlab.com/systemz/aimpanel2/master/response"
	"net/http"
)

func Register(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var registerRequest request.RegisterRequest
	err := decoder.Decode(&registerRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		lib.MustEncode(json.NewEncoder(w),
			response.JsonError{ErrorCode: 1001, Message: "Invalid body."})
		return
	}

	if registerRequest.Password != registerRequest.PasswordRepeat {
		w.WriteHeader(http.StatusBadRequest)

		lib.MustEncode(json.NewEncoder(w),
			response.JsonError{ErrorCode: 1002, Message: "Passwords do not match."})
		return
	}

	if registerRequest.Email != registerRequest.EmailRepeat {
		w.WriteHeader(http.StatusBadRequest)

		lib.MustEncode(json.NewEncoder(w),
			response.JsonError{ErrorCode: 1003, Message: "Emails do not match."})
		return
	}

	var count int64
	db.DB.Model(&model.User{}).Where("username = ?", registerRequest.Username).Count(&count)
	if count > 0 {
		w.WriteHeader(http.StatusBadRequest)

		lib.MustEncode(json.NewEncoder(w),
			response.JsonError{ErrorCode: 1004, Message: "User with this username already exist."})
		return
	}

	var user model.User
	user.Username = registerRequest.Username
	user.Email = registerRequest.Email
	user.PasswordHash = user.HashPassword(registerRequest.Password)

	err = db.DB.Save(&user).Error
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		lib.MustEncode(json.NewEncoder(w),
			response.JsonError{ErrorCode: 1005, Message: "Something went wrong."})
		return
	}

	token, err := user.GenerateJWT()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		lib.MustEncode(json.NewEncoder(w),
			response.JsonError{ErrorCode: 1006, Message: "Something went wrong."})
		return
	}

	lib.MustEncode(json.NewEncoder(w), response.TokenResponse{Token: token})
}

func Login(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var loginRequest request.LoginRequest
	err := decoder.Decode(&loginRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		lib.MustEncode(json.NewEncoder(w),
			response.JsonError{ErrorCode: 1007, Message: "Invalid body."})
		return
	}

	var user model.User
	db.DB.Where("username = ?", loginRequest.Username).Find(&user)

	if user.CheckPassword(loginRequest.Password) {
		token, err := user.GenerateJWT()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)

			lib.MustEncode(json.NewEncoder(w),
				response.JsonError{ErrorCode: 1008, Message: "Something went wrong."})
			return
		}

		lib.MustEncode(json.NewEncoder(w), response.TokenResponse{Token: token})
	} else {
		w.WriteHeader(http.StatusBadRequest)

		lib.MustEncode(json.NewEncoder(w),
			response.JsonError{ErrorCode: 1009, Message: "Wrong username or password."})
		return
	}
}
