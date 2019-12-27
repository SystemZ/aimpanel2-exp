package handler

import (
	"encoding/json"
	"github.com/gorilla/context"
	"gitlab.com/systemz/aimpanel2/lib"
	"gitlab.com/systemz/aimpanel2/lib/ecode"
	"gitlab.com/systemz/aimpanel2/master/model"
	"gitlab.com/systemz/aimpanel2/master/response"
	"net/http"
)

//swagger:parameters User changePassword
type UserChangePasswordReq struct {
	Password          string `json:"password"`
	NewPassword       string `json:"new_password"`
	NewPasswordRepeat string `json:"new_password_repeat"`
}

//swagger:parameters User changeEmail
type UserChangeEmailReq struct {
	Email          string `json:"email"`
	NewEmail       string `json:"new_email"`
	NewEmailRepeat string `json:"new_email_repeat"`
}

func UserChangePassword(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user").(model.User)

	decoder := json.NewDecoder(r.Body)
	var changePasswordReq UserChangePasswordReq
	err := decoder.Decode(&changePasswordReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		lib.MustEncode(json.NewEncoder(w),
			JsonError{ErrorCode: ecode.JsonDecode})
		return
	}

	if user.CheckPassword(changePasswordReq.Password) {
		if changePasswordReq.NewPassword == changePasswordReq.NewPasswordRepeat {
			user.PasswordHash = user.HashPassword(changePasswordReq.NewPassword)
			model.DB.Save(&user)

			w.WriteHeader(http.StatusOK)
			return
		} else {
			w.WriteHeader(http.StatusBadRequest)

			lib.MustEncode(json.NewEncoder(w),
				JsonError{ErrorCode: ecode.PasswordsDoNotMatch})
			return
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)

		lib.MustEncode(json.NewEncoder(w),
			JsonError{ErrorCode: ecode.WrongPassword})
		return
	}
}

func UserChangeEmail(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user").(model.User)

	decoder := json.NewDecoder(r.Body)
	var changeEmailReq UserChangeEmailReq
	err := decoder.Decode(&changeEmailReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		lib.MustEncode(json.NewEncoder(w),
			JsonError{ErrorCode: ecode.JsonDecode})
		return
	}

	if user.Email == changeEmailReq.Email {
		if changeEmailReq.NewEmail == changeEmailReq.NewEmailRepeat {
			user.Email = changeEmailReq.NewEmail
			model.DB.Save(&user)

			w.WriteHeader(http.StatusOK)
			return
		} else {
			w.WriteHeader(http.StatusBadRequest)

			lib.MustEncode(json.NewEncoder(w),
				JsonError{ErrorCode: ecode.EmailsDoNotMatch})
			return
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)

		lib.MustEncode(json.NewEncoder(w),
			JsonError{ErrorCode: ecode.WrongEmail})
		return
	}
}

// @Summary Profile Info
// @Tags User
// @Description Show currently logged in user details
// @Accept json
// @Produce json
// @Success 200 {object} response.UserProfile
// @Failure 400 {object} JsonError
// @Router /me [get]
func UserProfile(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user").(model.User)

	userProfile := response.UserProfileResponse{
		User: response.UserProfile{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		},
	}

	lib.MustEncode(json.NewEncoder(w), userProfile)
}
