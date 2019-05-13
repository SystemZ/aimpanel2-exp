package request

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
