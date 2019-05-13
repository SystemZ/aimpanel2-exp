package handler

import (
	"encoding/json"
	"github.com/gorilla/context"
	"gitlab.com/systemz/aimpanel2/lib"
	"gitlab.com/systemz/aimpanel2/master/db"
	"gitlab.com/systemz/aimpanel2/master/model"
	"gitlab.com/systemz/aimpanel2/master/request"
	"gitlab.com/systemz/aimpanel2/master/response"
	"net/http"
)

func ChangePassword(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user").(model.User)

	decoder := json.NewDecoder(r.Body)
	var changePasswordReq request.UserChangePasswordReq
	err := decoder.Decode(&changePasswordReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		lib.MustEncode(json.NewEncoder(w),
			response.JsonError{ErrorCode: 2001})
		return
	}

	if user.CheckPassword(changePasswordReq.Password) {
		if changePasswordReq.NewPassword == changePasswordReq.NewPasswordRepeat {
			user.PasswordHash = user.HashPassword(changePasswordReq.NewPassword)
			db.DB.Save(&user)

			w.WriteHeader(http.StatusOK)
			return
		} else {
			w.WriteHeader(http.StatusBadRequest)

			lib.MustEncode(json.NewEncoder(w),
				response.JsonError{ErrorCode: 2002})
			return
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)

		lib.MustEncode(json.NewEncoder(w),
			response.JsonError{ErrorCode: 2003})
		return
	}
}

func ChangeEmail(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user").(model.User)

	decoder := json.NewDecoder(r.Body)
	var changeEmailReq request.UserChangeEmailReq
	err := decoder.Decode(&changeEmailReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		lib.MustEncode(json.NewEncoder(w),
			response.JsonError{ErrorCode: 2004})
		return
	}

	if user.Email == changeEmailReq.Email {
		if changeEmailReq.NewEmail == changeEmailReq.NewEmailRepeat {
			user.Email = changeEmailReq.NewEmail
			db.DB.Save(&user)

			w.WriteHeader(http.StatusOK)
			return
		} else {
			w.WriteHeader(http.StatusBadRequest)

			lib.MustEncode(json.NewEncoder(w),
				response.JsonError{ErrorCode: 2005})
			return
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)

		lib.MustEncode(json.NewEncoder(w),
			response.JsonError{ErrorCode: 2006})
		return
	}
}

func Profile(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user").(model.User)
	lib.MustEncode(json.NewEncoder(w), user)
}
