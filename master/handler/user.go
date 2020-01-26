package handler

import (
	"encoding/json"
	"github.com/gorilla/context"
	"gitlab.com/systemz/aimpanel2/lib"
	"gitlab.com/systemz/aimpanel2/lib/ecode"
	"gitlab.com/systemz/aimpanel2/lib/request"
	"gitlab.com/systemz/aimpanel2/master/model"
	"gitlab.com/systemz/aimpanel2/master/response"
	userService "gitlab.com/systemz/aimpanel2/master/service/user"
	"net/http"
)

func UserChangePassword(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user").(model.User)

	data := &request.UserChangePasswordReq{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		lib.MustEncode(json.NewEncoder(w),
			JsonError{ErrorCode: ecode.JsonDecode})
		return
	}

	errCode := userService.ChangePassword(data, &user)
	if errCode != ecode.NoError {
		w.WriteHeader(http.StatusBadRequest)
		lib.MustEncode(json.NewEncoder(w), JsonError{ErrorCode: errCode})
		return
	}

	w.WriteHeader(http.StatusOK)
}

func UserChangeEmail(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user").(model.User)

	data := &request.UserChangeEmailReq{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		lib.MustEncode(json.NewEncoder(w),
			JsonError{ErrorCode: ecode.JsonDecode})
		return
	}

	errCode := userService.ChangeEmail(data, &user)
	if errCode != ecode.NoError {
		w.WriteHeader(http.StatusBadRequest)
		lib.MustEncode(json.NewEncoder(w), JsonError{ErrorCode: errCode})
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary Profile Info
// @Tags User
// @Description Show currently logged in user details
// @Accept json
// @Produce json
// @Success 200 {object} response.UserProfile
// @Failure 400 {object} JsonError
// @Router /me [get]
// @Security ApiKey
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
