package model

import (
	"context"
	"fmt"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/master/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var (
	DB       *mongo.Database
	DBOnline bool
)

func InitDB() *mongo.Database {
	clientOptions := options.Client()
	clientOptions.ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:%s/?authSource=admin", config.DB_USERNAME, config.DB_PASSWORD, config.DB_HOST, config.DB_PORT))

	if config.DEV_MODE {
		clientOptions.SetDirect(true)
	}

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		logrus.Error(err.Error())
		panic("Failed to connect to database")
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		logrus.Panic("Ping to db failed")
	}

	logrus.Info("Connection to database seems OK!")
	DBOnline = true

	db := client.Database(config.DB_NAME)

	go DBPing()

	return db
}

func DBPing() {
	for {
		<-time.After(5 * time.Second)

		err := DB.Client().Ping(context.TODO(), nil)
		if err != nil {
			DBOnline = false
		} else {
			DBOnline = true
		}
	}
}
