package wrapper

import (
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/lib/task"
	"gitlab.com/systemz/aimpanel2/slave/config"
	"gitlab.com/systemz/aimpanel2/slave/model"
)

func Start(gameServerID string) {
	logrus.Info("starting wrapper version " + config.GIT_COMMIT)

	output := make(chan string)
	input := make(chan string)

	p := &Process{
		Output:       output,
		Input:        input,
		GameServerID: gameServerID,
	}

	redisStarted := make(chan bool, 1)
	go p.RedisListener(redisStarted)
	<-redisStarted

	taskMsg := task.Message{
		TaskId:       task.GAME_STARTED,
		GameServerID: gameServerID,
	}
	model.SendTask(config.REDIS_PUB_SUB_AGENT_CH, taskMsg)

	taskMsg = task.Message{
		TaskId:       task.GAME_METRICS_FREQUENCY,
		GameServerID: gameServerID,
	}
	model.SendTask(config.REDIS_PUB_SUB_AGENT_CH, taskMsg)

	select {}
}
