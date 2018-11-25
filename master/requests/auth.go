package requests

//swagger:parameters Authentication register
type RegisterRequest struct {
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
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
