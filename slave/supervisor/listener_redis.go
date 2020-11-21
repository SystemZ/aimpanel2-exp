package supervisor

import (
	"context"
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/lib/task"
	"gitlab.com/systemz/aimpanel2/slave/config"
	"gitlab.com/systemz/aimpanel2/slave/model"
	"gitlab.com/systemz/aimpanel2/slave/tasks"
)

func listenerRedis(done chan bool) {
	// start connection to redis
	//model.InitRedis()

	// subscribe tasks
	// https://godoc.org/github.com/go-redis/redis#PubSub
	pubsub := model.Redis.Subscribe(context.TODO(), config.REDIS_PUB_SUB_SUPERVISOR_CH)

	// Wait for confirmation that subscription is created before publishing anything.
	_, err := pubsub.Receive(context.TODO())
	if err != nil {
		// FIXME don't panic on redis pub/sub error
		panic(err)
	}
	defer pubsub.Close()

	// Go channel which receives messages.
	ch := pubsub.Channel()

	done <- true

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
	case task.SUPERVISOR_CLEAN_FILES_GS:
		logrus.Info("supervisor got SUPERVISOR_CLEAN_FILES_GS")
		tasks.GsCleanFiles(taskMsg.GameServerID)

		taskMsg := task.Message{
			TaskId:       task.SUPERVISOR_CLEAN_FILES_GS,
			GameServerID: taskMsg.GameServerID,
		}
		model.SendTask(config.REDIS_PUB_SUB_AGENT_CH, taskMsg)
	default:
		logrus.Warningf("Unhandled task %v!", taskMsg.TaskId)
	}
}
