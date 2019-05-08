package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/master/config"
	"gitlab.com/systemz/aimpanel2/master/model"
)

var (
	DB *gorm.DB
)

func Setup() *gorm.DB {
	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True", config.DB_USERNAME, config.DB_PASSWORD, config.DB_HOST, config.DB_PORT, config.DB_NAME))
	if err != nil {
		panic("Failed to connect to database")
	}

	err = db.DB().Ping()
	if err != nil {
		logrus.Panic("Ping to db failed")
	}

	//https://github.com/go-sql-driver/mysql/issues/257
	db.DB().SetMaxIdleConns(0)

	db.LogMode(true)

	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Host{})

	db.AutoMigrate(&model.GameServer{})
	db.AutoMigrate(&model.GameServerLog{})

	db.AutoMigrate(&model.Game{})
	db.AutoMigrate(&model.GameFile{})
	db.AutoMigrate(&model.GameCommand{})
	db.AutoMigrate(&model.GameVersion{})

	db.AutoMigrate(&model.MetricHost{})
	db.AutoMigrate(&model.MetricGameServer{})

	logrus.Info("Connection to database seems OK!")

	return db
}
