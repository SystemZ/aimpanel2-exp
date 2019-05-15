package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/master/config"
)

var (
	DB *gorm.DB
)

func InitMysql() *gorm.DB {
	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True", config.DB_USERNAME, config.DB_PASSWORD, config.DB_HOST, config.DB_PORT, config.DB_NAME))
	if err != nil {
		logrus.Error(err.Error())
		panic("Failed to connect to database")
	}

	err = db.DB().Ping()
	if err != nil {
		logrus.Panic("Ping to db failed")
	}

	//https://github.com/go-sql-driver/mysql/issues/257
	db.DB().SetMaxIdleConns(0)

	db.LogMode(true)

	db.AutoMigrate(&User{})
	db.AutoMigrate(&Host{})

	db.AutoMigrate(&GameServer{})
	db.AutoMigrate(&GameServerLog{})

	db.AutoMigrate(&Game{})
	db.AutoMigrate(&GameFile{})
	db.AutoMigrate(&GameCommand{})
	db.AutoMigrate(&GameVersion{})

	db.AutoMigrate(&MetricHost{})
	db.AutoMigrate(&MetricGameServer{})

	logrus.Info("Connection to database seems OK!")

	return db
}
