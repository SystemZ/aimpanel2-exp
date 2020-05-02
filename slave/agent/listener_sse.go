package agent

import (
	"github.com/r3labs/sse"
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/lib"
	"gitlab.com/systemz/aimpanel2/lib/task"
	"gitlab.com/systemz/aimpanel2/slave/config"
	"gitlab.com/systemz/aimpanel2/slave/tasks"
	"strconv"
)

func listenerSse(done chan bool) {
	client := sse.NewClient(config.API_URL + "/v1/events/" + config.HOST_TOKEN)
	client.Headers = map[string]string{
		"Authorization": "Bearer " + config.API_TOKEN,
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
		taskId, _ := strconv.Atoi(string(msg.Event))

		switch taskId {
		case task.WRAPPER_START:
			// executed by agent
			logrus.Info("agent wrapper start requested")
			tasks.StartWrapper(taskMsg)
			logrus.Info("agent wrapper start finished")
		case task.GAME_INSTALL:
			// executed by agent
			logrus.Info("agent gs install requested")
			tasks.GsInstall(taskMsg)
			logrus.Info("agent gs install finished")
		case task.AGENT_METRICS_FREQUENCY:
			// executed by agent
			logrus.Infof("agent metrics freq %v sec requested", taskMsg.MetricFrequency)
			metricsFrequency = taskMsg.MetricFrequency
			go metrics()
			logrus.Info("agent metrics freq finished")
		case task.SLAVE_UPDATE:
			// executed by agent, it's closing whole agent process at the end
			logrus.Info("agent update requested")
			tasks.SelfUpdate(taskMsg)
		case task.AGENT_REMOVE_GS:
			// executed by agent
			logrus.Info("agent gs remove requested")
			tasks.GsRemove(taskMsg)
			logrus.Info("agent gs remove finished")
		case task.GAME_STOP_SIGTERM:
			// relayed to wrapper
			logrus.Info("agent sigterm requested")
			tasks.GsStop(taskMsg.GameServerID)
			logrus.Info("agent sigterm finished")
		case task.GAME_STOP_SIGKILL: // relayed to wrapper
			// relayed to wrapper
			logrus.Info("agent sigkill requested")
			tasks.GsKill(taskMsg.GameServerID)
			logrus.Info("agent sigkill finished")
		case task.GAME_COMMAND:
			// relayed to wrapper
			logrus.Info("agent gs cmd requested")
			tasks.GsCmd(taskMsg.GameServerID, taskMsg.Body)
			logrus.Info("agent gs cmd finished")
		case task.GAME_FILE_LIST:
			// executed by agent
			logrus.Info("agent gs file list requested")
			tasks.GsFileList(taskMsg.GameServerID)
			logrus.Info("agent gs file list finished")
		default:
			logrus.Info("Unknown task")
		}
	}
}
