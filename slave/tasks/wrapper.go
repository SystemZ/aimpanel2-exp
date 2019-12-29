package tasks

import (
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/lib/task"
	"gitlab.com/systemz/aimpanel2/slave/config"
	"gitlab.com/systemz/aimpanel2/slave/model"
)

// tasks below will be eventually finished by wrapper

func GsStart() {

}

func GsStop(gsId string) {
	taskMsg := task.Message{
		// FIXME other task IDs for user CLI actions
		TaskId:       task.GAME_STOP_SIGTERM,
		GameServerID: gsId,
	}
	taskMsgStr, err := taskMsg.Serialize()
	if err != nil {
		logrus.Errorf("preparing msg failed: %v", err)
		return
	}
	res, err := model.Redis.Publish(config.REDIS_PUB_SUB_CH, taskMsgStr).Result()
	if err != nil {
		logrus.Errorf("sending msg failed: %v", err)
	}
	logrus.Infof("Task sent to %v processes", res)
}

func GsKill(gsId string) {
	taskMsg := task.Message{
		// FIXME other task IDs for user CLI actions
		TaskId:       task.GAME_STOP_SIGKILL,
		GameServerID: gsId,
	}
	taskMsgStr, err := taskMsg.Serialize()
	if err != nil {
		logrus.Errorf("preparing msg failed: %v", err)
		return
	}
	res, err := model.Redis.Publish(config.REDIS_PUB_SUB_CH, taskMsgStr).Result()
	if err != nil {
		logrus.Errorf("sending msg failed: %v", err)
	}
	logrus.Infof("Task sent to %v processes", res)
}
