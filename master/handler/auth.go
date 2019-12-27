package handler

import (
	"encoding/json"
	"gitlab.com/systemz/aimpanel2/lib"
	"gitlab.com/systemz/aimpanel2/lib/ecode"
	"gitlab.com/systemz/aimpanel2/master/model"
	"gitlab.com/systemz/aimpanel2/master/response"
	"net/http"
)

type AuthRegisterRequest struct {
	Username       string `json:"username"`
	Password       string `json:"password"`
	PasswordRepeat string `json:"password_repeat"`
	Email          string `json:"email"`
	EmailRepeat    string `json:"email_repeat"`
}

type AuthLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Register(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var registerRequest AuthRegisterRequest
	err := decoder.Decode(&registerRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		lib.MustEncode(json.NewEncoder(w),
			JsonError{ErrorCode: ecode.JsonDecode})
		return
	}

	if registerRequest.Password != registerRequest.PasswordRepeat {
		w.WriteHeader(http.StatusBadRequest)

		lib.MustEncode(json.NewEncoder(w),
			JsonError{ErrorCode: ecode.WrongPassword})
		return
	}

	if registerRequest.Email != registerRequest.EmailRepeat {
		w.WriteHeader(http.StatusBadRequest)

		lib.MustEncode(json.NewEncoder(w),
			JsonError{ErrorCode: ecode.WrongEmail})
		return
	}

	var count int64
	model.DB.Model(&model.User{}).Where("username = ?", registerRequest.Username).Count(&count)
	if count > 0 {
		w.WriteHeader(http.StatusBadRequest)

		lib.MustEncode(json.NewEncoder(w),
			JsonError{ErrorCode: ecode.DuplicateUsername})
		return
	}

	var user model.User
	user.Username = registerRequest.Username
	user.Email = registerRequest.Email
	user.PasswordHash = user.HashPassword(registerRequest.Password)

	err = model.DB.Save(&user).Error
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		lib.MustEncode(json.NewEncoder(w),
			JsonError{ErrorCode: ecode.DbSave})
		return
	}

	token, err := user.GenerateJWT()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		lib.MustEncode(json.NewEncoder(w),
			JsonError{ErrorCode: ecode.JwtGenerate})
		return
	}

	//Create group
	group := &model.Group{
		Name: "USER-" + user.ID.String(),
	}
	model.DB.Save(group)

	//Add user to group
	groupUser := &model.GroupUser{
		GroupId: group.ID,
		UserId:  user.ID,
	}
	// FIXME error handling
	model.DB.Save(groupUser)
	model.CreatePermissionsForNewUser(group.ID)

	lib.MustEncode(json.NewEncoder(w), response.Token{Token: token})
}

func Login(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var loginRequest AuthLoginRequest
	err := decoder.Decode(&loginRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		lib.MustEncode(json.NewEncoder(w),
			JsonError{ErrorCode: ecode.JsonDecode})
		return
	}

	var user model.User
	model.DB.Where("username = ?", loginRequest.Username).Find(&user)

	// TODO maybe IsPasswordOk would be more semantic?
	if user.CheckPassword(loginRequest.Password) {
		token, err := user.GenerateJWT()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)

			lib.MustEncode(json.NewEncoder(w),
				JsonError{ErrorCode: ecode.JwtGenerate})
			return
		}

		lib.MustEncode(json.NewEncoder(w), response.Token{Token: token})
	} else {
		w.WriteHeader(http.StatusBadRequest)

		lib.MustEncode(json.NewEncoder(w),
			JsonError{ErrorCode: ecode.WrongUsernameOrPassword})
		return
	}
}
