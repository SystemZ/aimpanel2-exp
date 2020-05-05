package agent

import (
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/lib/ahttp"
	"gitlab.com/systemz/aimpanel2/lib/task"
	"gitlab.com/systemz/aimpanel2/slave/config"
	"gitlab.com/systemz/aimpanel2/slave/cron"
	"gitlab.com/systemz/aimpanel2/slave/model"
	"gitlab.com/systemz/aimpanel2/slave/tasks"
)

func listenerRedis(done chan bool) {
	// start connection to redis
	model.InitRedis()

	// subscribe tasks
	// https://godoc.org/github.com/go-redis/redis#PubSub
	pubsub := model.Redis.Subscribe(config.REDIS_PUB_SUB_AGENT_CH)

	// Wait for confirmation that subscription is created before publishing anything.
	_, err := pubsub.Receive()
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
	case task.GAME_STARTED:
		logrus.Info("Agent got " + taskMsg.TaskId.String())

		val, err := model.GetGsStart(taskMsg.GameServerID)
		if val == 1 {
			model.SetGsStart(taskMsg.GameServerID, 2)
			tasks.GsStartGame(taskMsg)
			model.DelGsStart(taskMsg.GameServerID)

			_, err = ahttp.SendTaskData(config.API_URL+"/v1/events/"+config.HOST_TOKEN, config.API_TOKEN, taskMsg)
			if err != nil {
				logrus.Error(err)
			}
		}

	case task.GAME_SHUTDOWN:
		logrus.Info("Agent got " + taskMsg.TaskId.String())

		_, err = ahttp.SendTaskData(config.API_URL+"/v1/events/"+config.HOST_TOKEN, config.API_TOKEN, taskMsg)
		if err != nil {
			logrus.Error(err)
		}

		//Start wrapper if gs is restarting
		val, _ := model.GetGsRestart(taskMsg.GameServerID)
		if val == 1 {
			logrus.Info("no to cyk")
			model.SetGsRestart(taskMsg.GameServerID, 2)
			tasks.StartWrapper(taskMsg)
			model.DelGsRestart(taskMsg.GameServerID)
		}

	case task.GAME_SERVER_LOG:
		logrus.Info("Agent got " + taskMsg.TaskId.String())

		_, err = ahttp.SendTaskData(config.API_URL+"/v1/events/"+config.HOST_TOKEN, config.API_TOKEN, taskMsg)
		if err != nil {
			logrus.Error(err)
		}
	case task.GAME_METRICS:
		logrus.Info("Agent got " + taskMsg.TaskId.String())

		_, err = ahttp.SendTaskData(config.API_URL+"/v1/events/"+config.HOST_TOKEN, config.API_TOKEN, taskMsg)
		if err != nil {
			logrus.Error(err)
		}
	case task.GAME_METRICS_FREQUENCY:
		logrus.Info("Agent got " + taskMsg.TaskId.String())

		_, err = ahttp.SendTaskData(config.API_URL+"/v1/events/"+config.HOST_TOKEN, config.API_TOKEN, taskMsg)
		if err != nil {
			logrus.Error(err)
		}

	case task.AGENT_BACKUP_GS:
		logrus.Info("agent got GAME_MAKE_BACKUP msg")
		go tasks.GsBackup(taskMsg.GameServerID)
	case task.AGENT_SHUTDOWN:
		logrus.Info("agent got AGENT_SHUTDOWN msg")
		cron.Stop()
		tasks.AgentShutdown()
	default:
		logrus.Warningf("Unhandled task %v!", taskMsg.TaskId)
	}
}
