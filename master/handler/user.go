package handler

import (
	"encoding/json"
	"gitlab.com/systemz/aimpanel2/master/db"
	"gitlab.com/systemz/aimpanel2/master/model"
	"gitlab.com/systemz/aimpanel2/master/request"
	"gitlab.com/systemz/aimpanel2/master/response"
	"log"
	"net/http"
)

func ChangePassword(w http.ResponseWriter, r *http.Request) {
	// swagger:route POST /user/change_password User changePassword
	//
	// Changes authenticated user password
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
	decoder := json.NewDecoder(r.Body)
	var changePasswordReq request.ChangePasswordReq
	err := decoder.Decode(&changePasswordReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response.JsonError{ErrorCode: 11, Message: "Invalid body."})
		return
	}

	var user model.User
	db.DB.Where("id = ?", r.Header.Get("uid")).First(&user)

	if user.CheckPassword(changePasswordReq.Password) {
		if changePasswordReq.NewPassword == changePasswordReq.NewPasswordRepeat {
			user.PasswordHash = user.HashPassword(changePasswordReq.NewPassword)
			db.DB.Save(&user)

			w.WriteHeader(http.StatusOK)
			return
		} else {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response.JsonError{ErrorCode: 12, Message: "Passwords do not match."})
			return
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response.JsonError{ErrorCode: 13, Message: "Current password is wrong."})
		return
	}
}

func ChangeEmail(w http.ResponseWriter, r *http.Request) {
	// swagger:route POST /user/change_email User changeEmail
	//
	// Changes authenticated user email
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
	decoder := json.NewDecoder(r.Body)
	var changeEmailReq request.ChangeEmailReq
	err := decoder.Decode(&changeEmailReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response.JsonError{ErrorCode: 14, Message: "Invalid body."})
		return
	}

	userId := r.Header.Get("uid")
	var user model.User
	db.DB.Where("id = ?", userId).First(&user)

	log.Println(user.Email)
	log.Println(changeEmailReq.Email)
	if user.Email == changeEmailReq.Email {
		if changeEmailReq.NewEmail == changeEmailReq.NewEmailRepeat {
			user.Email = changeEmailReq.NewEmail
			db.DB.Save(&user)

			w.WriteHeader(http.StatusOK)
			return
		} else {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response.JsonError{ErrorCode: 15, Message: "Emails do not match."})
			return
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response.JsonError{ErrorCode: 16, Message: "Current email is wrong."})
		return
	}
}

func Profile(w http.ResponseWriter, r *http.Request) {
	var user model.User
	db.DB.Where("id = ?", r.Header.Get("uid")).First(&user)

	json.NewEncoder(w).Encode(user)
}
