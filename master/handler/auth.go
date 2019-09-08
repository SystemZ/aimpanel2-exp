package handler

import (
	"encoding/json"
	"gitlab.com/systemz/aimpanel2/lib"
	"gitlab.com/systemz/aimpanel2/master/model"
	"net/http"
)

//swagger:response tokenResponse
type TokenResponse struct {
	Token string `json:"token"`
}

//swagger:parameters Authentication register
type AuthRegisterReq struct {
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

//swagger:parameters Authentication login
type AuthLoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// swagger:route POST /auth/register Auth Register
//
// Create new account
//
//Responses:
//	default: jsonError
//  400: jsonError
//	200:

func Register(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var registerRequest AuthRegisterReq
	err := decoder.Decode(&registerRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		lib.MustEncode(json.NewEncoder(w),
			JsonError{ErrorCode: 1001})
		return
	}

	if registerRequest.Password != registerRequest.PasswordRepeat {
		w.WriteHeader(http.StatusBadRequest)

		lib.MustEncode(json.NewEncoder(w),
			JsonError{ErrorCode: 1002})
		return
	}

	if registerRequest.Email != registerRequest.EmailRepeat {
		w.WriteHeader(http.StatusBadRequest)

		lib.MustEncode(json.NewEncoder(w),
			JsonError{ErrorCode: 1003})
		return
	}

	var count int64
	model.DB.Model(&model.User{}).Where("username = ?", registerRequest.Username).Count(&count)
	if count > 0 {
		w.WriteHeader(http.StatusBadRequest)

		lib.MustEncode(json.NewEncoder(w),
			JsonError{ErrorCode: 1004})
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
			JsonError{ErrorCode: 1005})
		return
	}

	token, err := user.GenerateJWT()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		lib.MustEncode(json.NewEncoder(w),
			JsonError{ErrorCode: 1006})
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
	model.DB.Save(groupUser)

	model.DB.Save(&model.Permission{
		Name:     "List hosts",
		Verb:     lib.GetVerbByName("GET"),
		GroupId:  group.ID,
		Endpoint: "/v1/host",
	})

	model.DB.Save(&model.Permission{
		Name:     "Create host",
		Verb:     lib.GetVerbByName("POST"),
		GroupId:  group.ID,
		Endpoint: "/v1/host",
	})

	model.DB.Save(&model.Permission{
		Name:     "List game servers by user id",
		Verb:     lib.GetVerbByName("GET"),
		GroupId:  group.ID,
		Endpoint: "/v1/host/my/server",
	})

	model.DB.Save(&model.Permission{
		Name:     "Change password",
		Verb:     lib.GetVerbByName("POST"),
		GroupId:  group.ID,
		Endpoint: "/v1/user/change_password",
	})

	model.DB.Save(&model.Permission{
		Name:     "Change email",
		Verb:     lib.GetVerbByName("POST"),
		GroupId:  group.ID,
		Endpoint: "/v1/user/change_email",
	})

	model.DB.Save(&model.Permission{
		Name:     "User profile",
		Verb:     lib.GetVerbByName("GET"),
		GroupId:  group.ID,
		Endpoint: "/v1/user/profile",
	})

	model.DB.Save(&model.Permission{
		Name:     "List games",
		Verb:     lib.GetVerbByName("GET"),
		GroupId:  group.ID,
		Endpoint: "/v1/game",
	})

	model.DB.Save(&model.Permission{
		Name:     "List game versions",
		Verb:     lib.GetVerbByName("GET"),
		GroupId:  group.ID,
		Endpoint: "/v1/game/version",
	})

	lib.MustEncode(json.NewEncoder(w), TokenResponse{Token: token})
}

// swagger:route POST /auth/login Auth Login
//
// Get session token
//
//Responses:
//	default: jsonError
//  400: jsonError
//	200:

func Login(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var loginRequest AuthLoginReq
	err := decoder.Decode(&loginRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		lib.MustEncode(json.NewEncoder(w),
			JsonError{ErrorCode: 1007})
		return
	}

	var user model.User
	model.DB.Where("username = ?", loginRequest.Username).Find(&user)

	if user.CheckPassword(loginRequest.Password) {
		token, err := user.GenerateJWT()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)

			lib.MustEncode(json.NewEncoder(w),
				JsonError{ErrorCode: 1008})
			return
		}

		lib.MustEncode(json.NewEncoder(w), TokenResponse{Token: token})
	} else {
		w.WriteHeader(http.StatusBadRequest)

		lib.MustEncode(json.NewEncoder(w),
			JsonError{ErrorCode: 1009})
		return
	}
}
