package model

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/lib/game"
	"gitlab.com/systemz/aimpanel2/lib/task"
	"gitlab.com/systemz/aimpanel2/slave/config"
	"time"
)

var (
	Redis *redis.Client
)

func InitRedis() {
	Redis = redis.NewClient(&redis.Options{
		Network:  "unix",
		Addr:     config.REDIS_HOST,
		Username: config.REDIS_USERNAME,
		Password: config.REDIS_PASSWORD,
	})
	_, err := Redis.Ping(context.Background()).Result()
	if err != nil {
		logrus.Error(err.Error())
		logrus.Panic("Ping to Redis failed")
	}
	logrus.Info("Connected to DB")
}

func SendTask(channel string, taskMsg task.Message) bool {
	taskMsgStr, _ := taskMsg.Serialize()

	res, err := Redis.Publish(context.TODO(), channel, taskMsgStr).Result()
	if err != nil {
		logrus.Errorf("sending %v failed: %v", taskMsg.TaskId, err)
		return false
	}

	logrus.Infof("Task %v sent to %v processes", taskMsg.TaskId, res)

	return true
}

func SetGsRestart(gsId string, state int) {
	Redis.Set(context.TODO(), "gs_restart_id_"+gsId, state, 24*time.Hour)
}

func GetGsRestart(gsId string) (int64, error) {
	return Redis.Get(context.TODO(), "gs_restart_id_"+gsId).Int64()
}

func DelGsRestart(gsId string) {
	Redis.Del(context.TODO(), "gs_restart_id_"+gsId)
}

func SetGsGame(gsId string, game *game.Game) {
	gameStr, err := json.Marshal(game)
	if err != nil {
		logrus.Errorf("error while saving game: %v", err)
		return
	}

	Redis.Set(context.TODO(), "gs_"+gsId+"_game", string(gameStr), 24*time.Hour)
}

func GetGsGame(gsId string) (game.Game, error) {
	var g game.Game
	gameStr, err := Redis.Get(context.TODO(), "gs_"+gsId+"_game").Result()
	err = json.Unmarshal([]byte(gameStr), &g)
	return g, err
}

func SetGsStart(gsId string, state int) {
	Redis.Set(context.TODO(), "gs_start_id_"+gsId, state, 24*time.Hour)
}

func GetGsStart(gsId string) (int64, error) {
	return Redis.Get(context.TODO(), "gs_start_id_"+gsId).Int64()
}

func DelGsStart(gsId string) {
	Redis.Del(context.TODO(), "gs_start_id_"+gsId)
}

func SetGsRunning(gsId string, running int) {
	Redis.Set(context.TODO(), "gs_running_id_"+gsId, running, 0)
}

func GetGsRunning(gsId string) (int64, error) {
	return Redis.Get(context.TODO(), "gs_running_id_"+gsId).Int64()
}

func SetHwId(hwId string) {
	Redis.Set(context.TODO(), "hw_id", hwId, 0)
}

func GetHwId() string {
	val, err := Redis.Get(context.TODO(), "hw_id").Result()
	if err != nil {
		return ""
	}
	return val
}

func SetHostToken(hostToken string) {
	Redis.Set(context.TODO(), "host_token", hostToken, 0)
}

func GetHostToken() string {
	val, err := Redis.Get(context.TODO(), "host_token").Result()
	if err != nil {
		logrus.Warn(err)
		return ""
	}
	return val
}
