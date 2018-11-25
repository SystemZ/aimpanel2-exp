package requests

type ChangePasswordReq struct {
	Password          string `json:"password"`
	NewPassword       string `json:"new_password"`
	NewPasswordRepeat string `json:"new_password_repeat"`
}

type ChangeEmailReq struct {
	Email          string `json:"email"`
	NewEmail       string `json:"new_email"`
	NewEmailRepeat string `json:"new_email_repeat"`
}
