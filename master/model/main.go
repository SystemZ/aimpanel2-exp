package model

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/master/config"

	"github.com/bwmarrin/snowflake"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Redis     *redis.Client
	DB        *mongo.Database
	Snowflake *snowflake.Node
)

func InitDB() *mongo.Database {
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:%s/?authSource=admin", config.DB_USERNAME, config.DB_PASSWORD, config.DB_HOST, config.DB_PORT))

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

func InitSnowflake() *snowflake.Node {
	node, err := snowflake.NewNode(config.NODE_ID)
	if err != nil {
		logrus.Error(err.Error())
		panic("Failed to set snowflake node")
	}

	return node
}
