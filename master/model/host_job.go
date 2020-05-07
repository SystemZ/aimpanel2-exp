package model

import (
	"context"
	"gitlab.com/systemz/aimpanel2/lib/task"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HostJob struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty" example:"1238206236281802752"`

	//User assigned name
	Name string `bson:"name" json:"name" example:"Restart server"`

	// Host ID
	//
	// required: true
	HostId primitive.ObjectID `bson:"host_id" json:"host_id" example:"100112233-4455-6677-8899-aabbccddeeff"`

	CronExpression string `bson:"cron_expression" json:"cron_expression" example:"5 4 * * *"`

	TaskMessage task.Message `bson:"task_message" json:"task_message"`
}

func (h *HostJob) GetCollectionName() string {
	return hostJobCollection
}

func (h *HostJob) GetID() primitive.ObjectID {
	return h.ID
}

func (h *HostJob) SetID(id primitive.ObjectID) {
	h.ID = id
}

func GetHostJobsByHostId(hostId primitive.ObjectID) ([]HostJob, error) {
	var hostJobs []HostJob

	cur, err := DB.Collection(hostJobCollection).Find(context.TODO(),
		bson.D{{"host_id", hostId}})
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var hostJob HostJob
		if err := cur.Decode(&hostJob); err != nil {
			return nil, err
		}
		hostJobs = append(hostJobs, hostJob)
	}

	return hostJobs, nil
}

func GetHostJobById(id primitive.ObjectID) (*HostJob, error) {
	var hostJob HostJob

	err := DB.Collection(hostJobCollection).FindOne(context.TODO(), bson.D{{"_id", id}}).Decode(&hostJob)
	if err != nil {
		return nil, err
	}

	return &hostJob, nil
}
