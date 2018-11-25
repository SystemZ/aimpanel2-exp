package handlers

import (
	"encoding/json"
	"gitlab.com/systemz/aimpanel2/master/db"
	"gitlab.com/systemz/aimpanel2/master/models"
	"gitlab.com/systemz/aimpanel2/master/requests"
	"gitlab.com/systemz/aimpanel2/master/responses"
	"log"
	"net/http"
)

func ChangePassword(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var changePasswordReq requests.ChangePasswordReq
	err := decoder.Decode(&changePasswordReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(responses.JsonError{ErrorCode: 11, Message: "Invalid body."})
		return
	}

	var user models.User
	db.DB.Where("id = ?", r.Header.Get("uid")).First(&user)

	if user.CheckPassword(changePasswordReq.Password) {
		if changePasswordReq.NewPassword == changePasswordReq.NewPasswordRepeat {
			user.PasswordHash = user.HashPassword(changePasswordReq.NewPassword)
			db.DB.Save(&user)

			w.WriteHeader(http.StatusOK)
			return
		} else {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(responses.JsonError{ErrorCode: 12, Message: "Passwords do not match."})
			return
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(responses.JsonError{ErrorCode: 13, Message: "Current password is wrong."})
		return
	}
}

func ChangeEmail(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var changeEmailReq requests.ChangeEmailReq
	err := decoder.Decode(&changeEmailReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(responses.JsonError{ErrorCode: 14, Message: "Invalid body."})
		return
	}

	userId := r.Header.Get("uid")
	var user models.User
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
			json.NewEncoder(w).Encode(responses.JsonError{ErrorCode: 15, Message: "Emails do not match."})
			return
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(responses.JsonError{ErrorCode: 16, Message: "Current email is wrong."})
		return
	}
}

func Profile(w http.ResponseWriter, r *http.Request) {
	var user models.User
	db.DB.Where("id = ?", r.Header.Get("uid")).First(&user)

	json.NewEncoder(w).Encode(user)
}
