package gameserver

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/alexandrevicenzi/go-sse"
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/lib"
	"gitlab.com/systemz/aimpanel2/lib/ecode"
	"gitlab.com/systemz/aimpanel2/lib/task"
	"gitlab.com/systemz/aimpanel2/master/events"
	"gitlab.com/systemz/aimpanel2/master/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Data(hostToken string, taskMsg *task.Message) error {
	switch taskMsg.TaskId {
	case task.GAME_STARTED:
		logrus.Infof("Got %v", taskMsg.TaskId)

	case task.GAME_SERVER_LOG:
		logrus.Infof("Got %v", taskMsg.TaskId)
		err := Log(hostToken, taskMsg)
		if err != nil {
			return err
		}
	case task.GAME_SHUTDOWN:
		logrus.Infof("Got %v", taskMsg.TaskId)

	case task.GAME_METRICS_FREQUENCY:
		logrus.Infof("Got %v", taskMsg.TaskId)
		err := GameMetricsFrequency(hostToken, taskMsg)
		if err != nil {
			return err
		}
	case task.GAME_METRICS:
		logrus.Infof("Got %v", taskMsg.TaskId)
		err := Metrics(*taskMsg)
		if err != nil {
			return err
		}

	case task.AGENT_FILE_LIST_GS:
		logrus.Infof("Got %v", taskMsg.TaskId)
		err := model.GsFilesPublish(model.Redis, taskMsg.GameServerID, taskMsg.Files)
		if err != nil {
			logrus.Error(err)
		}
	default:
		logrus.Infof("Unhandled task %v", taskMsg.TaskId)
	}

	return nil
}

func Log(hostToken string, taskMsg *task.Message) error {
	var gsLog model.GameServerLog
	oid, _ := primitive.ObjectIDFromHex(taskMsg.GameServerID)
	gsLog.GameServerId = oid

	if len(taskMsg.Stdout) > 0 {
		gsLog.Log = taskMsg.Stdout
		gsLog.Type = model.STDOUT
	}

	if len(taskMsg.Stderr) > 0 {
		gsLog.Log = taskMsg.Stderr
		gsLog.Type = model.STDERR
	}

	host, err := model.GetHostByToken(hostToken)
	if err != nil {
		return err
	}

	//TODO: Handle this with mongo events collection
	events.SSE.SendMessage(fmt.Sprintf("/v1/host/%s/server/%s/console",
		host.ID.Hex(),
		gsLog.GameServerId.Hex()),
		sse.SimpleMessage(base64.StdEncoding.EncodeToString([]byte(gsLog.Log))))

	err = model.Put(&gsLog)
	if err != nil {
		logrus.Warn(err)
	}

	return nil
}

func GameMetricsFrequency(hostToken string, taskMsg *task.Message) error {
	gameServerId := taskMsg.GameServerID

	host, err := model.GetHostByToken(hostToken)
	if err != nil {
		return &lib.Error{ErrorCode: ecode.HostNotFound}
	}

	oid, err := primitive.ObjectIDFromHex(taskMsg.GameServerID)
	if err != nil {
		return err
	}

	gs, err := model.GetGameServerById(oid)
	if err != nil {
		return err
	}

	if gs == nil {
		return errors.New("game server not found")
	}

	taskMessage := task.Message{
		TaskId:          task.GAME_METRICS_FREQUENCY,
		GameServerID:    gameServerId,
		MetricFrequency: gs.MetricFrequency,
	}

	err = model.SendEvent(host.ID, taskMessage)
	if err != nil {
		return &lib.Error{ErrorCode: ecode.DbSave}
	}

	return nil
}
