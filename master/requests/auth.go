package requests

type RegisterRequest struct {
	Username       string `json:"username"`
	Password       string `json:"password"`
	PasswordRepeat string `json:"password_repeat"`
	Email          string `json:"email"`
	EmailRepeat    string `json:"email_repeat"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
