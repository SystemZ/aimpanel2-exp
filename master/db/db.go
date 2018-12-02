package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gitlab.com/systemz/aimpanel2/master/models"
	"log"
	"os"
)

var (
	DB       *gorm.DB
	hostname = os.Getenv("DB_HOSTNAME")
	name     = os.Getenv("DB_NAME")
	username = os.Getenv("DB_USERNAME")
	password = os.Getenv("DB_PASSWORD")
)

func SetupDatabase() *gorm.DB {
	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True", username, password, hostname, name))
	if err != nil {
		panic("Failed to connect to database")
	}

	//https://github.com/go-sql-driver/mysql/issues/257
	db.DB().SetMaxIdleConns(0)

	db.LogMode(true)

	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Host{})

	db.AutoMigrate(&models.Game{})
	db.AutoMigrate(&models.GameEngine{})

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
