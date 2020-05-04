package agent

import (
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/lib"
	"gitlab.com/systemz/aimpanel2/lib/ahttp"
	"gitlab.com/systemz/aimpanel2/lib/response"
	"gitlab.com/systemz/aimpanel2/lib/task"
	"gitlab.com/systemz/aimpanel2/slave/config"
	"gitlab.com/systemz/aimpanel2/slave/cron"
	"gitlab.com/systemz/aimpanel2/slave/model"
	"gitlab.com/systemz/aimpanel2/slave/tasks"
	"time"
)

func Start(hostToken string) {
	model.InitRedis()
	cron.InitCron()
	ahttp.HttpClient = ahttp.InitHttpClient()

	logrus.Info("Starting Agent Version." + config.GIT_COMMIT)
	config.HOST_TOKEN = hostToken

	var token response.Token
	_, err := ahttp.Get(config.API_URL+"/v1/host/auth/"+config.HOST_TOKEN, &token)
	if err != nil {
		lib.FailOnError(err, "Failed to get host token")
	}
	config.API_TOKEN = token.Token

	sseStarted := make(chan bool, 1)
	redisStarted := make(chan bool, 1)
	go listenerSse(sseStarted)
	go listenerRedis(redisStarted)

	<-sseStarted
	<-redisStarted

	logrus.Info("Send AGENT_STARTED")
	taskMsg := task.Message{
		TaskId: task.AGENT_STARTED,
	}
	//TODO: do something with status code
	_, err = ahttp.SendTaskData(config.API_URL+"/v1/events/"+config.HOST_TOKEN, config.API_TOKEN, taskMsg)
	if err != nil {
		logrus.Error(err)
	}

	logrus.Info("Send AGENT_METRICS_FREQUENCY")
	taskMsg = task.Message{
		TaskId: task.AGENT_METRICS_FREQUENCY,
	}
	//TODO: do something with status code
	_, err = ahttp.SendTaskData(config.API_URL+"/v1/events/"+config.HOST_TOKEN, config.API_TOKEN, taskMsg)
	if err != nil {
		logrus.Error(err)
	}

	go func() {
		time.Sleep(time.Duration(lib.RandInt(200, 2000)) * time.Millisecond)

		logrus.Info("Send AGENT_GET_JOBS")
		taskMsg = task.Message{
			TaskId: task.AGENT_GET_JOBS,
		}

		//TODO: do something with status code
		_, err = ahttp.SendTaskData(config.API_URL+"/v1/events/"+config.HOST_TOKEN, config.API_TOKEN, taskMsg)
		if err != nil {
			logrus.Error(err)
		}
	}()

	tasks.AgentSendOSInfo()

	select {}
}
