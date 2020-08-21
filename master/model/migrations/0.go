package migrations

import (
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/master/model"
)

func MigrateUp() (err error) {
	configExists, err := model.IsConfigDocPresent()
	if err != nil {
		return err
	}
	if !configExists {
		err := Migration1Up()
		if err != nil {
			logrus.Error(err)
		}
	}
	return
}
