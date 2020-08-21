package auth

import (
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/lib/ecode"
	"gitlab.com/systemz/aimpanel2/lib/request"
	"gitlab.com/systemz/aimpanel2/master/model"
)

func Register(data *request.AuthRegister) (string, int) {
	if data.Password != data.PasswordRepeat {
		return "", ecode.WrongPassword
	}

	if data.Email != data.EmailRepeat {
		return "", ecode.WrongEmail
	}

	if model.CheckIfUserExist(data.Username) {
		return "", ecode.DuplicateUsername
	}

	var user model.User
	user.Username = data.Username
	user.Email = data.Email
	user.PasswordHash = user.HashPassword(data.Password)

	err := model.Put(&user)
	if err != nil {
		logrus.Error(err)
		return "", ecode.DbSave
	}

	token, err := user.GenerateJWT()
	if err != nil {
		return "", ecode.JwtGenerate
	}

	/*
		err = model.CreatePermissionsForNewUser(user.ID)
		if err != nil {
			logrus.Error(err)
			return "", ecode.DbSave
		}
	*/

	return token, ecode.NoError
}

func Login(data *request.AuthLogin) (string, int) {
	user, err := model.GetUserByUsername(data.Username)
	if err != nil {
		logrus.Error(err)
		return "", ecode.DbError
	}

	if user.IsPasswordOk(data.Password) {
		token, err := user.GenerateJWT()
		if err != nil {
			return "", ecode.JwtGenerate
		}

		return token, ecode.NoError
	} else {
		return "", ecode.WrongUsernameOrPassword
	}
}
