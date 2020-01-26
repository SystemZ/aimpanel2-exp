package request

type HostCreate struct {
	// User assigned name
	Name string `json:"name" example:"My Great Linux server"`

	// User assigned ip
	Ip string `json:"ip" example:"192.51.100.128"`
}

type GameServerCreate struct {
	//User assigned name
	Name string `json:"name" example:"Ultra MC Server"`

	//Selected game id
	GameId uint `json:"game_id" example:"1"`

	//Selected game version
	GameVersion string `json:"game_version" example:"1.14.2"`
}

type GameServerStop struct {
	Type uint `json:"type" example:"1"`
}

type GameServerSendCommand struct {
	Command string `json:"command" example:"say Hello!"`
}

//swagger:parameters User changePassword
type UserChangePassword struct {
	Password          string `json:"password"`
	NewPassword       string `json:"new_password"`
	NewPasswordRepeat string `json:"new_password_repeat"`
}

//swagger:parameters User changeEmail
type UserChangeEmail struct {
	Email          string `json:"email"`
	NewEmail       string `json:"new_email"`
	NewEmailRepeat string `json:"new_email_repeat"`
}

type AuthRegister struct {
	Username       string `json:"username"`
	Password       string `json:"password"`
	PasswordRepeat string `json:"password_repeat"`
	Email          string `json:"email"`
	EmailRepeat    string `json:"email_repeat"`
}

type AuthLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
