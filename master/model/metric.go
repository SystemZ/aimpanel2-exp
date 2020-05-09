package model

import (
	"context"
	"gitlab.com/systemz/aimpanel2/lib/metric"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

//metric := &model.MetricHost{
//	HostId:    host.ID,
//	CpuUsage:  taskMsg.CpuUsage,
//	RamFree:   taskMsg.RamFree,
//	RamTotal:  taskMsg.RamTotal,
//	DiskFree:  taskMsg.DiskFree,
//	DiskUsed:  taskMsg.DiskUsed,
//	DiskTotal: taskMsg.DiskTotal,
//	User:      taskMsg.User,
//	System:    taskMsg.System,
//	Idle:      taskMsg.Idle,
//	Nice:      taskMsg.Nice,
//	Iowait:    taskMsg.Iowait,
//	Irq:       taskMsg.Irq,
//	Softirq:   taskMsg.Softirq,
//	Steal:     taskMsg.Steal,
//	Guest:     taskMsg.Guest,
//	GuestNice: taskMsg.GuestNice,
//}
//err = model.Put(metric)
//if err != nil {
//	return err
//}

const (
	GameServerMetric = iota
	HostMetric
)

type Metric struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty" example:"1238206236281802752"`
	Type     uint8              `bson:"type" json:"type"`
	RID      primitive.ObjectID `bson:"r_id" json:"r_id"`
	Metric   metric.Id          `bson:"metric" json:"metric"`
	NSamples int                `bson:"nsamples" json:"nsamples"`
	Day      time.Time          `bson:"day" json:"day"`
	First    int64              `bson:"first" json:"first"`
	Last     int64              `bson:"last" json:"last"`
	Samples  []MetricData       `bson:"samples" json:"samples"`
}

func (m *Metric) GetCollectionName() string {
	return metricCollection
}

func (m *Metric) GetID() primitive.ObjectID {
	return m.ID
}

type MetricData struct {
	Val  int   `bson:"val" json:"val"`
	Time int64 `bson:"time" json:"time"`
}

func PutMetric(metricType uint8, rid primitive.ObjectID, metric metric.Id, val int, metricTime int64) error {
	opts := options.Update()
	opts.SetUpsert(true)

	now := time.Now()
	day := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	metricData := MetricData{
		Val:  val,
		Time: metricTime,
	}

	filter := bson.D{
		{Key: "type", Value: metricType},
		{Key: "r_id", Value: rid},
		{Key: "metric", Value: metric},
	}

	update := bson.D{
		{Key: "$push", Value: bson.D{{Key: "samples", Value: metricData}}},
		{Key: "$min", Value: bson.D{{Key: "first", Value: metricData.Time}}},
		{Key: "$max", Value: bson.D{{Key: "last", Value: metricData.Time}}},
		{Key: "$inc", Value: bson.D{{Key: "nsamples", Value: 1}}},
		{Key: "$set", Value: bson.D{{Key: "day", Value: day}}},
	}

	_, err := DB.Collection(metricCollection).UpdateOne(context.TODO(), filter, update, opts)

	return err
}
