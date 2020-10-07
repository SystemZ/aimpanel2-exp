package migrations

import (
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/master/model"
)

func MigrateUp() (err error) {
	// 0 -> 1
	configExists, err := model.IsConfigDocPresent()
	if err != nil {
		return err
	}
	if !configExists {
		logrus.Info("Applying migration ID 1")
		err = Migration1Up()
		if err != nil {
			return
		}
	}

	// read config from DB
	config, err := model.GetConfig()
	if err != nil {
		return err
	}

	// 1 -> 2
	if config.MigrateVersion < 2 {
		logrus.Info("Applying migration ID 2")
		err = Migration2Up()
		if err != nil {
			return
		}
		config.MigrateVersion = 2
		model.Update(&config)
	}
	// 2 -> 3
	if config.MigrateVersion < 3 {
		logrus.Info("Applying migration ID 3")
		err = Migration3Up()
		if err != nil {
			return
		}
		config.MigrateVersion = 3
		model.Update(&config)
	}
	// 3 -> 4
	if config.MigrateVersion < 4 {
		logrus.Info("Applying migration ID 4")
		err = Migration4Up()
		if err != nil {
			return
		}
		config.MigrateVersion = 4
		model.Update(&config)
	}
	// 4 -> 5
	if config.MigrateVersion < 5 {
		logrus.Info("Applying migration ID 5")
		err = Migration5Up()
		if err != nil {
			return
		}
		config.MigrateVersion = 5
		model.Update(&config)
	}

	// end
	return
}
