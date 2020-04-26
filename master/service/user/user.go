package user

import (
	"gitlab.com/systemz/aimpanel2/lib/ecode"
	"gitlab.com/systemz/aimpanel2/lib/request"
	"gitlab.com/systemz/aimpanel2/master/model"
)

func ChangePassword(data *request.UserChangePassword, user *model.User) int {
	if user.IsPasswordOk(data.Password) {
		if data.NewPassword == data.NewPasswordRepeat {
			user.PasswordHash = user.HashPassword(data.NewPassword)
			err := user.Update(&user)
			if err != nil {
				return ecode.DbSave
			}

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
			err := user.Update(&user)
			if err != nil {
				return ecode.DbSave
			}

			return ecode.NoError
		} else {
			return ecode.EmailsDoNotMatch
		}
	} else {
		return ecode.WrongEmail
	}
}
