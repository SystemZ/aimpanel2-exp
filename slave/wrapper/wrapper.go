package wrapper

import (
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/lib"
	"gitlab.com/systemz/aimpanel2/lib/task"
	"gitlab.com/systemz/aimpanel2/slave/config"
)

func Start(gameServerID string) {
	logrus.Info("Starting Wrapper Version. " + config.GIT_COMMIT)

	output := make(chan string)
	input := make(chan string)

	p := &Process{
		Output:       output,
		Input:        input,
		GameServerID: gameServerID,
	}

	go p.Rpc()

	logrus.Info("Send WRAPPER_STARTED")
	taskMsg := task.Message{
		TaskId:       task.WRAPPER_STARTED,
		GameServerID: gameServerID,
	}

	jsonStr, err := taskMsg.Serialize()
	if err != nil {
		logrus.Error(err)
	}
	//TODO: do something with status code
	_, err = lib.SendTaskData(config.API_URL+"/v1/events/"+config.HOST_TOKEN+"/"+gameServerID, config.API_TOKEN, jsonStr)
	if err != nil {
		logrus.Error(err)
	}

	logrus.Info("Send WRAPPER_METRICS_FREQUENCY")
	taskMsg = task.Message{
		TaskId:       task.WRAPPER_METRICS_FREQUENCY,
		GameServerID: gameServerID,
	}

	jsonStr, err = taskMsg.Serialize()
	if err != nil {
		logrus.Error(err)
	}
	//TODO: do something with status code
	_, err = lib.SendTaskData(config.API_URL+"/v1/events/"+config.HOST_TOKEN+"/"+gameServerID, config.API_TOKEN, jsonStr)
	if err != nil {
		logrus.Error(err)
	}

	select {}
}
