package agent

import (
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/lib/task"
	"gitlab.com/systemz/aimpanel2/slave/config"
	"gitlab.com/systemz/aimpanel2/slave/model"
)

func listenerRedis() {
	// start connection to redis
	model.InitRedis()

	// subscribe tasks
	// https://godoc.org/github.com/go-redis/redis#PubSub
	pubsub := model.Redis.Subscribe(config.REDIS_PUB_SUB_CH)

	// Wait for confirmation that subscription is created before publishing anything.
	_, err := pubsub.Receive()
	if err != nil {
		// FIXME don't panic on redis pub/sub error
		panic(err)
	}
	defer pubsub.Close()

	// Go channel which receives messages.
	ch := pubsub.Channel()

	// Consume messages.
	for msg := range ch {
		redisTaskHandler(msg.Channel, msg.Payload)
	}
}

func redisTaskHandler(taskCh string, taskBody string) {
	taskMsg := task.Message{}
	err := taskMsg.Deserialize(taskBody)
	if err != nil {
		logrus.Error(err)
	}

	switch taskMsg.TaskId {

	case task.GAME_STOP_SIGTERM:
		logrus.Info("agent got GAME_STOP_SIGTERM msg")
	case task.GAME_STOP_SIGKILL:
		logrus.Info("agent got GAME_STOP_SIGKILL msg")
	case task.GAME_COMMAND:
		logrus.Info("agent got GAME_COMMAND msg")
	default:
		logrus.Warning("Unhandled task!")
	}
}
