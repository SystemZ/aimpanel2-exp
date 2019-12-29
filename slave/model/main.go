package model

import (
	"github.com/go-redis/redis/v7"
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/slave/config"
)

var (
	Redis *redis.Client
)

func InitRedis() {
	Redis = redis.NewClient(&redis.Options{
		Network: "unix",
		Addr:    config.REDIS_HOST,
	})
	_, err := Redis.Ping().Result()
	if err != nil {
		logrus.Error(err.Error())
		logrus.Panic("Ping to Redis failed")
	}
	logrus.Info("Connected to DB")
}
