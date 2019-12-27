package handler

import (
	"encoding/json"
	"gitlab.com/systemz/aimpanel2/lib"
	"gitlab.com/systemz/aimpanel2/lib/ecode"
	"gitlab.com/systemz/aimpanel2/master/model"
	"net/http"
)

//swagger:model registerRequest
type AuthRegisterRequest struct {
	// required: true
	Username string `json:"username"`
	// required: true
	Password string `json:"password"`
	// required: true
	PasswordRepeat string `json:"password_repeat"`
	// required: true
	Email string `json:"email"`
	// required: true
	EmailRepeat string `json:"email_repeat"`
}

//swagger:parameters Auth Login
type AuthLoginRequest struct {
	// User name
	//
	// in: body
	// required: true
	Username string `json:"username"`

	// User password
	//
	// in: body
	// required: true
	Password string `json:"password"`
}

//swagger:model Token
type Token struct {
	// JWT token
	Token string `json:"token"`
}

// swagger:route POST /auth/register Auth Register
//
// Create new account
//
//Responses:
//	default: jsonError
//  400: jsonError
//	200: Token

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

	lib.MustEncode(json.NewEncoder(w), Token{Token: token})
}

// swagger:route POST /auth/login Auth Login
//
// Get session token
//
//Responses:
//	default: jsonError
//  400: jsonError
//	200: Token

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

		lib.MustEncode(json.NewEncoder(w), Token{Token: token})
	} else {
		w.WriteHeader(http.StatusBadRequest)

		lib.MustEncode(json.NewEncoder(w),
			JsonError{ErrorCode: ecode.WrongUsernameOrPassword})
		return
	}
}
