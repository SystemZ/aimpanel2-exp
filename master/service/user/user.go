package user

import (
	"gitlab.com/systemz/aimpanel2/lib/ecode"
	"gitlab.com/systemz/aimpanel2/lib/request"
	"gitlab.com/systemz/aimpanel2/master/model"
)

func ChangePassword(data *request.UserChangePassword, user *model.User) int {
	if user.CheckPassword(data.Password) {
		if data.NewPassword == data.NewPasswordRepeat {
			user.PasswordHash = user.HashPassword(data.NewPassword)
			model.DB.Save(&user)

			return ecode.NoError
		} else {
			return ecode.PasswordsDoNotMatch
		}
	} else {
		return ecode.WrongPassword
	}
}

func ChangeEmail(data *request.UserChangeEmail, user *model.User) int {
	if user.Email == data.Email {
		if data.NewEmail == data.NewEmailRepeat {
			user.Email = data.NewEmail
			model.DB.Save(&user)

			return ecode.NoError
		} else {
			return ecode.EmailsDoNotMatch
		}
	} else {
		return ecode.WrongEmail
	}
}
