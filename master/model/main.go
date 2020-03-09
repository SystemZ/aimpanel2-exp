package model

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/master/config"

	_ "github.com/go-kivik/couchdb/v3" // The CouchDB driver
	"github.com/go-kivik/kivik/v3"
)

var (
	DB      *gorm.DB
	Redis   *redis.Client
	CouchDB *kivik.DB
)

func InitCouchDb() *kivik.DB {
	client, err := kivik.New("couch", fmt.Sprintf("http://%s:%s@%s:%s/", config.DB_USERNAME, config.DB_PASSWORD, config.DB_HOST, config.DB_PORT))
	if err != nil {
		logrus.Error(err.Error())
		panic("Failed to connect to database")
	}

	_, err = client.Ping(context.TODO())
	if err != nil {
		logrus.Panic("Ping to db failed")
	}

	db := client.DB(context.TODO(), config.DB_NAME)
	logrus.Info("Connection to database seems OK!")

	return db
}

func InitMysql() *gorm.DB {
	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True", config.DB_USERNAME, config.DB_PASSWORD, config.DB_HOST, "3306", config.DB_NAME))
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
	db.LogMode(config.DEV_MODE)

	db.AutoMigrate(&User{})
	db.AutoMigrate(&Host{})

	db.AutoMigrate(&GameServer{})
	db.AutoMigrate(&GameServerLog{})

	db.AutoMigrate(&GameFile{})

	db.AutoMigrate(&MetricHost{})
	db.AutoMigrate(&MetricGameServer{})

	db.AutoMigrate(&Group{})
	db.AutoMigrate(&GroupUser{})

	db.AutoMigrate(&Permission{})

	db.AutoMigrate(&Event{})

	logrus.Info("Connection to database seems OK!")

	return db
}

func InitRedis() {
	if len(config.REDIS_PASSWORD) > 1 {
		Redis = redis.NewClient(&redis.Options{
			Addr:     config.REDIS_HOST + ":6379",
			Password: config.REDIS_PASSWORD,
		})
	} else {
		Redis = redis.NewClient(&redis.Options{
			Addr: config.REDIS_HOST + ":6379",
		})
	}

	_, err := Redis.Ping().Result()
	if err != nil {
		logrus.Error(err.Error())
		logrus.Panic("Ping to Redis failed")
	}

	logrus.Info("Connection to Redis seems OK!")
}
