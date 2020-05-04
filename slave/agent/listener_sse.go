package agent

import (
	"github.com/r3labs/sse"
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/lib"
	"gitlab.com/systemz/aimpanel2/lib/ahttp"
	"gitlab.com/systemz/aimpanel2/lib/task"
	"gitlab.com/systemz/aimpanel2/slave/config"
	"gitlab.com/systemz/aimpanel2/slave/cron"
	"gitlab.com/systemz/aimpanel2/slave/tasks"
	"net/http"
)

func listenerSse(done chan bool) {
	client := sse.NewClient(config.API_URL + "/v1/events/" + config.HOST_TOKEN)
	client.Headers = map[string]string{
		"Authorization": "Bearer " + config.API_TOKEN,
	}
	client.Connection.Transport = &http.Transport{
		DialTLSContext: ahttp.VerifyPinTLSContext,
	}

	events := make(chan *sse.Event)
	err := client.SubscribeChan("", events)
	if err != nil {
		lib.FailOnError(err, "Can't connect to event channel")
	}

	logrus.Info("Subscribed to SSE")
	done <- true

	for msg := range events {
		logrus.Info(msg.ID)
		logrus.Info(string(msg.Data))
		logrus.Info(string(msg.Event))

		taskMsg := task.Message{}
		err := taskMsg.Deserialize(string(msg.Data))
		if err != nil {
			logrus.Error(err)
		}

		if taskMsg.TaskId == task.AGENT_GET_JOBS {
			go cron.AddJobs(taskMsg.Jobs)

			continue
		}

		tasks.ProcessTask(taskMsg)
	}
}
