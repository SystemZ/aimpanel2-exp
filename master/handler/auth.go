package handler

import (
	"encoding/json"
	"gitlab.com/systemz/aimpanel2/master/db"
	"gitlab.com/systemz/aimpanel2/master/model"
	"gitlab.com/systemz/aimpanel2/master/request"
	"gitlab.com/systemz/aimpanel2/master/response"
	"net/http"
)

func Register(w http.ResponseWriter, r *http.Request) {
	// swagger:route POST /auth/register Authentication register
	//
	// Registers new account
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
	//	200: tokenResponse
	decoder := json.NewDecoder(r.Body)
	var registerRequest request.RegisterRequest
	err := decoder.Decode(&registerRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response.JsonError{ErrorCode: 1, Message: "Invalid body."})
		return
	}

	if registerRequest.Password != registerRequest.PasswordRepeat {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response.JsonError{ErrorCode: 2, Message: "Passwords do not match."})
		return
	}

	if registerRequest.Email != registerRequest.EmailRepeat {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response.JsonError{ErrorCode: 3, Message: "Emails do not match."})
		return
	}

	var count int64
	db.DB.Model(&model.User{}).Where("username = ?", registerRequest.Username).Count(&count)
	if count > 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response.JsonError{ErrorCode: 4, Message: "User with this username already exist."})
		return
	}

	var user model.User
	user.Username = registerRequest.Username
	user.Email = registerRequest.Email
	user.PasswordHash = user.HashPassword(registerRequest.Password)

	err = db.DB.Save(&user).Error
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response.JsonError{ErrorCode: 5, Message: "Something went wrong."})
		return
	}

	token, err := user.GenerateJWT()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response.JsonError{ErrorCode: 6, Message: "Something went wrong."})
		return
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(response.TokenResponse{Token: token})
}

func Login(w http.ResponseWriter, r *http.Request) {
	// swagger:route POST /auth/login Authentication login
	//
	// Authenticates the user
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
	//	200: tokenResponse
	decoder := json.NewDecoder(r.Body)
	var loginRequest request.LoginRequest
	err := decoder.Decode(&loginRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response.JsonError{ErrorCode: 7, Message: "Invalid body."})
		return
	}

	var user model.User
	db.DB.Where("username = ?", loginRequest.Username).Find(&user)

	if user.CheckPassword(loginRequest.Password) {
		token, err := user.GenerateJWT()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response.JsonError{ErrorCode: 8, Message: "Something went wrong."})
			return
		}

		w.WriteHeader(200)
		json.NewEncoder(w).Encode(response.TokenResponse{Token: token})
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response.JsonError{ErrorCode: 9, Message: "Wrong password."})
		return
	}
}
