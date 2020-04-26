package model

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/master/config"

	_ "github.com/go-kivik/couchdb/v3" // The CouchDB driver
	"github.com/go-kivik/kivik/v3"

	"github.com/bwmarrin/snowflake"
)

var (
	Redis     *redis.Client
	DB        *kivik.DB
	Snowflake *snowflake.Node
)

func InitDB() *kivik.DB {
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
