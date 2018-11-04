package handlers

import (
	"encoding/json"
	"gitlab.com/systemz/aimpanel2/master/db"
	"gitlab.com/systemz/aimpanel2/master/models"
	"gitlab.com/systemz/aimpanel2/master/requests"
	"gitlab.com/systemz/aimpanel2/master/responses"
	"net/http"
)

func Register(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var registerRequest requests.RegisterRequest
	err := decoder.Decode(&registerRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(responses.JsonError{ErrorCode: 1, Message: "Invalid body."})
		return
	}

	if registerRequest.Password != registerRequest.PasswordRepeat {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(responses.JsonError{ErrorCode: 2, Message: "Passwords do not match."})
		return
	}

	if registerRequest.Email != registerRequest.EmailRepeat {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(responses.JsonError{ErrorCode: 3, Message: "Emails do not match."})
		return
	}

	var count int64
	db.DB.Model(&models.User{}).Where("username = ?", registerRequest.Username).Count(&count)
	if count > 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(responses.JsonError{ErrorCode: 4, Message: "User with this username already exist."})
		return
	}

	var user models.User
	user.Username = registerRequest.Username
	user.Email = registerRequest.Email
	user.PasswordHash = user.HashPassword(registerRequest.Password)

	err = db.DB.Save(&user).Error
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(responses.JsonError{ErrorCode: 5, Message: "Something went wrong."})
		return
	}

	token, err := user.GenerateJWT()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(responses.JsonError{ErrorCode: 6, Message: "Something went wrong."})
		return
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(responses.TokenResponse{Token: token})
}
