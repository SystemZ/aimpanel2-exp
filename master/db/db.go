package db

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/jinzhu/gorm"

	"os"
	"fmt"
)

var (
	DB       *gorm.DB
	hostname = os.Getenv("DB_HOSTNAME")
	name     = os.Getenv("DB_NAME")
	username = os.Getenv("DB_USERNAME")
	password = os.Getenv("DB_PASSWORD")
)

func SetupDatabase() *gorm.DB {
	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s%%tcp(%s)/%s?charset=utf8&parseTime=True", username, password, hostname, name))
	if err != nil {
		panic("Failed to connect to database")
	}

	//https://github.com/go-sql-driver/mysql/issues/257
	db.DB().SetMaxIdleConns(0)

	return db
}
