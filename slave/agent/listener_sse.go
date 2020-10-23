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
	connected := false

	ahttpHost := ahttp.Hosts[ahttp.CurrentHost]
	client := sse.NewClient(ahttpHost + "/v1/events/" + config.HOST_TOKEN)
	client.Headers = map[string]string{
		"Authorization": config.HW_ID,
	}
	client.Connection.Transport = &http.Transport{
		DialTLSContext: ahttp.VerifyPinTLSContext,
	}
	client.ReconnectStrategy = NewUnlimitedRetry(true, "sse")
	client.OnDisconnect(func(c *sse.Client) {
		connected = false
		logrus.Warn("SSE disconnected")
	})

	events := make(chan *sse.Event)
	err := client.SubscribeChan("", events)
	if err != nil {
		lib.FailOnError(err, "Can't connect to event channel")
	}

	logrus.Info("Subscribed to SSE")
	connected = true
	done <- true

	for msg := range events {
		if !connected {
			logrus.Info("Reconnected to SSE successfully")
			connected = true
		}
		taskMsg := task.Message{}
		err := taskMsg.Deserialize(string(msg.Data))
		if err != nil {
			logrus.Error(err)
		}

		if taskMsg.TaskId == task.AGENT_GET_JOBS {
			go cron.AddJobs(taskMsg.Jobs)
			continue
		}

		go tasks.ProcessTask(taskMsg)
	}
}
