package agent

import (
	"github.com/r3labs/sse"
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/lib"
	"gitlab.com/systemz/aimpanel2/lib/ahttp"
	"gitlab.com/systemz/aimpanel2/lib/task"
	"gitlab.com/systemz/aimpanel2/slave/config"
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

		switch taskMsg.TaskId {
		case task.GAME_COMMAND, task.GAME_STOP_SIGKILL,
			task.GAME_STOP_SIGTERM, task.GAME_RESTART, task.GAME_METRICS_FREQUENCY:

			// executed by wrapper
			tasks.WrapperTaskHandler(taskMsg)
		case task.AGENT_START_GS, task.AGENT_INSTALL_GS,
			task.AGENT_BACKUP_GS, task.AGENT_UPDATE,
			task.AGENT_REMOVE_GS, task.AGENT_FILE_LIST_GS:

			// executed by agent
			tasks.AgentTaskHandler(taskMsg)

		case task.AGENT_METRICS_FREQUENCY:
			// executed by agent
			logrus.Infof("agent metrics freq %v sec requested", taskMsg.MetricFrequency)
			metricsFrequency = taskMsg.MetricFrequency
			go metrics()
			logrus.Info("agent metrics freq finished")
		default:
			logrus.Infof("Unknown task %v", taskMsg.TaskId)
		}
	}
}
