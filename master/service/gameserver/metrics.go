package gameserver

import (
	"gitlab.com/systemz/aimpanel2/lib/metric"
	"gitlab.com/systemz/aimpanel2/lib/task"
	"gitlab.com/systemz/aimpanel2/master/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func Metrics(taskMsg task.Message) error {
	oid, err := primitive.ObjectIDFromHex(taskMsg.GameServerID)
	if err != nil {
		return err
	}

	err = model.PutMetric(model.GameServerMetric, oid, metric.CpuUsage, taskMsg.CpuUsage, time.Now().Unix())
	if err != nil {
		return err
	}

	err = model.PutMetric(model.GameServerMetric, oid, metric.RamUsage, taskMsg.RamUsage, time.Now().Unix())
	if err != nil {
		return err
	}

	return nil
}
