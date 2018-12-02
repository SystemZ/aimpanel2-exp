package request

//swagger:parameters User changePassword
type ChangePasswordReq struct {
	Password          string `json:"password"`
	NewPassword       string `json:"new_password"`
	NewPasswordRepeat string `json:"new_password_repeat"`
}

//swagger:parameters User changeEmail
type ChangeEmailReq struct {
	Email          string `json:"email"`
	NewEmail       string `json:"new_email"`
	NewEmailRepeat string `json:"new_email_repeat"`
}
