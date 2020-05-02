package wrapper

import (
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/lib/ahttp"
	"gitlab.com/systemz/aimpanel2/lib/task"
	"gitlab.com/systemz/aimpanel2/slave/config"
)

func (p *Process) LogStdout(msg string) {
	taskMsg := task.Message{
		TaskId:       task.SERVER_LOG,
		GameServerID: p.GameServerID,
		Stdout:       msg,
	}

	jsonStr, err := taskMsg.Serialize()
	if err != nil {
		logrus.Error(err)
	}
	//TODO: do something with status code
	_, err = ahttp.SendTaskData(config.API_URL+"/v1/events/"+config.HOST_TOKEN+"/"+p.GameServerID, config.API_TOKEN, jsonStr)
	if err != nil {
		logrus.Error(err)
	}
}

func (p *Process) LogStderr(msg string) {
	taskMsg := task.Message{
		TaskId:       task.SERVER_LOG,
		GameServerID: p.GameServerID,
		Stderr:       msg,
	}

	jsonStr, err := taskMsg.Serialize()
	if err != nil {
		logrus.Error(err)
	}
	//TODO: do something with status code
	_, err = ahttp.SendTaskData(config.API_URL+"/v1/events/"+config.HOST_TOKEN+"/"+p.GameServerID, config.API_TOKEN, jsonStr)
	if err != nil {
		logrus.Error(err)
	}
}
