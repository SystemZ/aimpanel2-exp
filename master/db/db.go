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

	db.AutoMigrate(&model.Game{})
	db.AutoMigrate(&model.GameFile{})
	db.AutoMigrate(&model.GameCommand{})

	db.AutoMigrate(&model.GameServerLog{})

	logrus.Info("Connection to database seems OK!")

	//Create test user
	//user := models.User{Username: "test", Password: "test"}
	//db.Create(&user)

	//Find user
	//var user models.User
	//db.Where("ID = ?", "99a8f14e-1d5f-4de5-8812-2f23e5b1f446").First(&user)
	//log.Println(user)

	return db
}
