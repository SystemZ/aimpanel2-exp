package model

import (
	"context"
	"fmt"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/master/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	DB *mongo.Database
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

	db := client.Database(config.DB_NAME)

	return db
}
