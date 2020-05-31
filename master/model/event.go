package model

import (
	"context"
	"fmt"
	"github.com/alexandrevicenzi/go-sse"
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/lib"
	"gitlab.com/systemz/aimpanel2/lib/task"
	"gitlab.com/systemz/aimpanel2/master/events"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type Event struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty" example:"1238206236281802752"`

	HostId primitive.ObjectID `bson:"host_id"`

	TaskMessage task.Message `bson:"task_message" json:"task_message"`
}

func (e *Event) GetCollectionName() string {
	return eventCollection
}

func (e *Event) GetID() primitive.ObjectID {
	return e.ID
}

func (e *Event) SetID(id primitive.ObjectID) {
	e.ID = id
}

func SendEvent(hostId primitive.ObjectID, taskMsg task.Message) error {
	event := &Event{
		HostId:      hostId,
		TaskMessage: taskMsg,
	}
	err := Put(event)
	if err != nil {
		return err
	}

	return nil
}

func EventChanges() {
	pipeline := bson.D{
		{
			Key: "$match",
			Value: bson.D{
				{Key: "operationType", Value: "insert"},
			},
		},
	}

	changeStream, err := DB.Collection(eventCollection).Watch(context.TODO(), mongo.Pipeline{pipeline})
	if err != nil {
		log.Fatal(err)
	}
	defer changeStream.Close(context.TODO())

	for changeStream.Next(context.TODO()) {
		logrus.Info("EventChanges() Got event")
		changeDoc := struct {
			FullDocument Event `bson:"fullDocument"`
		}{}

		if err := changeStream.Decode(&changeDoc); err != nil {
			logrus.Fatal(err)
			continue
		}

		if changeDoc.FullDocument.HostId == primitive.NilObjectID {
			logrus.Info("EventChanges() HostId is empty. Ignoring.")
			continue
		}

		availableHosts := events.SSE.Channels()

		host, err := GetHostById(changeDoc.FullDocument.HostId)
		if err != nil {
			logrus.Warn("EventChanges() Could not find host. Ignoring.")
			continue
		}

		channelName := fmt.Sprintf("/v1/events/%s", host.Token)
		if !lib.StringInSlice(channelName, availableHosts) {
			logrus.Info("EventChanges() Event not for me. Ignoring.")
			continue
		}

		channel, _ := events.SSE.GetChannel(channelName)
		taskMsg := changeDoc.FullDocument.TaskMessage
		taskMsgStr, err := taskMsg.Serialize()
		if err != nil {
			logrus.Warn("EventChanges() Could not serialize task message. Ignoring.")
		}

		channel.SendMessage(sse.NewMessage("", taskMsgStr, taskMsg.TaskId.StringValue()))

		logrus.Infof("EventChanges() Task sent to host %v", host.ID.Hex())
	}

	if err := changeStream.Err(); err != nil {
		log.Fatal(err)
	}
}
