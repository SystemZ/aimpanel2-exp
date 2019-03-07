package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gitlab.com/systemz/aimpanel2/master/config"
	"gitlab.com/systemz/aimpanel2/master/model"
	"log"
)

var (
	DB *gorm.DB
)

func SetupDatabase() *gorm.DB {
	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True", config.DB_USERNAME, config.DB_PASSWORD, config.DB_HOST, config.DB_PORT, config.DB_NAME))
	if err != nil {
		panic("Failed to connect to database")
	}

	//https://github.com/go-sql-driver/mysql/issues/257
	db.DB().SetMaxIdleConns(0)

	db.LogMode(true)

	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Host{})

	db.AutoMigrate(&model.Game{})
	db.AutoMigrate(&model.GameEngine{})

	log.Println("Connected to the database")

	//Create test user
	//user := models.User{Username: "test", Password: "test"}
	//db.Create(&user)

	//Find user
	//var user models.User
	//db.Where("ID = ?", "99a8f14e-1d5f-4de5-8812-2f23e5b1f446").First(&user)
	//log.Println(user)

	return db
}
