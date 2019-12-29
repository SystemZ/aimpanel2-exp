package agent

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/lib"
	"gitlab.com/systemz/aimpanel2/lib/response"
	"gitlab.com/systemz/aimpanel2/lib/task"
	"gitlab.com/systemz/aimpanel2/slave/config"
	"gitlab.com/systemz/aimpanel2/slave/model"
	"net/http"
)

var (
	metricsFrequency int
)

func Start(hostToken string) {
	model.InitRedis()

	logrus.Info("Starting Agent Version." + config.GIT_COMMIT)
	config.HOST_TOKEN = hostToken

	resp, err := http.Get(config.API_URL + "/v1/host/auth/" + config.HOST_TOKEN)
	if err != nil {
		lib.FailOnError(err, "Failed to get host token")
	}

	var token response.Token
	err = json.NewDecoder(resp.Body).Decode(&token)
	if err != nil {
		lib.FailOnError(err, "Failed to decode credentials json")
	}
	config.API_TOKEN = token.Token

	go listenerSse()

	logrus.Info("Send AGENT_STARTED")
	taskMsg := task.Message{
		TaskId: task.AGENT_STARTED,
	}

	jsonStr, err := taskMsg.Serialize()
	if err != nil {
		logrus.Error(err)
	}
	//TODO: do something with status code
	_, err = lib.SendTaskData(config.API_URL+"/v1/events/"+config.HOST_TOKEN, config.API_TOKEN, jsonStr)
	if err != nil {
		logrus.Error(err)
	}

	logrus.Info("Send AGENT_METRICS_FREQUENCY")
	taskMsg = task.Message{
		TaskId: task.AGENT_METRICS_FREQUENCY,
	}

	jsonStr, err = taskMsg.Serialize()
	if err != nil {
		logrus.Error(err)
	}
	//TODO: do something with status code
	_, err = lib.SendTaskData(config.API_URL+"/v1/events/"+config.HOST_TOKEN, config.API_TOKEN, jsonStr)
	if err != nil {
		logrus.Error(err)
	}

	sendOSInfo()

	select {}
}
