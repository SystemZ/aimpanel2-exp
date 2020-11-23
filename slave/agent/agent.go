package agent

import (
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/lib"
	"gitlab.com/systemz/aimpanel2/lib/ahttp"
	"gitlab.com/systemz/aimpanel2/lib/task"
	"gitlab.com/systemz/aimpanel2/slave/config"
	"gitlab.com/systemz/aimpanel2/slave/cron"
	"gitlab.com/systemz/aimpanel2/slave/model"
	"gitlab.com/systemz/aimpanel2/slave/tasks"
	"time"
)

func Start() {
	//model.InitRedis()
	logrus.Info("starting agent version " + config.GIT_COMMIT)
	cron.InitCron()
	ahttp.SlaveVersion = config.GIT_COMMIT
	ahttp.HttpClient = ahttp.InitHttpClient()

	//recreate directories if not exist
	tasks.SelfHeal()

	//Init redis
	model.InitRedis()

	// helper for sending tasks to supervisor
	model.EmitterInit()

	//Migrate host token from env to redis
	if len(config.HOST_TOKEN) > 0 && len(model.GetHostToken()) == 0 {
		model.SetHostToken(config.HOST_TOKEN)
	}

	hostToken := model.GetHostToken()
	config.HOST_TOKEN = model.GetHostToken()

	//Get HW ID from Redis, If empty create new one
	hwId := model.GetHwId()
	if hwId == "" {
		u, err := uuid.NewV4()
		if err != nil {
			lib.FailOnError(err, "Failed to generate HW ID")
		}

		model.SetHwId(u.String())
		config.HW_ID = u.String()
	}
	config.HW_ID = hwId

	// task which needs API token but SSE and redis isn't necessary
	//

	sseStarted := make(chan bool, 1)
	redisStarted := make(chan bool, 1)
	// all tasks from master are handled here
	go listenerSse(sseStarted)
	// wrapper and cli handling
	go listenerRedis(redisStarted)

	<-sseStarted
	<-redisStarted
	// FIXME handle redis in offline mode without master

	// tasks which needs just redis connected
	//

	// tasks which needs just SSE connected
	//

	// tasks which needs both SSE and redis already connected
	logrus.Info("Send AGENT_STARTED")
	taskMsg := task.Message{
		TaskId: task.AGENT_STARTED,
	}
	//TODO: do something with status code
	_, err := ahttp.SendTaskData("/v1/events/"+hostToken, config.HW_ID, taskMsg)
	if err != nil {
		logrus.Error(err)
	}

	logrus.Info("Send AGENT_METRICS_FREQUENCY")
	taskMsg = task.Message{
		TaskId: task.AGENT_METRICS_FREQUENCY,
	}
	//TODO: do something with status code
	_, err = ahttp.SendTaskData("/v1/events/"+hostToken, config.HW_ID, taskMsg)
	if err != nil {
		logrus.Error(err)
	}

	// TODO request jobs after every SSE DC
	// current way can desync jobs after slave is disconnected prolonged time
	// we can clear some redis key value when SSE disconnects in client.OnDisconnect()
	// then do endless loop with sleep and ask master when this key is empty
	// set timestamp after successful retrieval to skip gathering again
	// do it at start of slave regardless of timestamp
	// AGENT_METRICS_FREQUENCY probably need this too
	// AGENT_RECONNECTED event would be useful
	go func() {
		time.Sleep(time.Duration(lib.RandInt(200, 2000)) * time.Millisecond)

		logrus.Info("Send AGENT_GET_JOBS")
		taskMsg = task.Message{
			TaskId: task.AGENT_GET_JOBS,
		}

		//TODO: do something with status code
		_, err = ahttp.SendTaskData("/v1/events/"+hostToken, config.HW_ID, taskMsg)
		if err != nil {
			logrus.Error(err)
		}
	}()

	tasks.AgentSendOSInfo()

	select {}
}
