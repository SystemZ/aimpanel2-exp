package model

import (
	"fmt"
	"github.com/alexandrevicenzi/go-sse"
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/lib"
	"gitlab.com/systemz/aimpanel2/lib/task"
	"gitlab.com/systemz/aimpanel2/master/events"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Event struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty" example:"1238206236281802752"`
	UserId       primitive.ObjectID `bson:"user_id"`
	HostId       primitive.ObjectID `bson:"host_id"`
	GameServerId primitive.ObjectID `bson:"gs_id"`
	TaskId       int                `bson:"task_id"`
	Value        string             `bson:"val"`
	ValuePrev    string             `bson:"val_prev"`
}

func (e *Event) GetCollectionName() string {
	return EventCollection
}

func (e *Event) GetID() primitive.ObjectID {
	return e.ID
}

func (e *Event) SetID(id primitive.ObjectID) {
	e.ID = id
}

func SaveAction(taskMsg task.Message, user User, hostId primitive.ObjectID, val string, valPrev string) error {
	gsId, _ := primitive.ObjectIDFromHex(taskMsg.GameServerID)
	event := &Event{
		UserId:       user.ID,
		TaskId:       int(taskMsg.TaskId),
		HostId:       hostId,
		GameServerId: gsId,
		Value:        val,
		ValuePrev:    valPrev,
	}
	// note additional details
	switch taskMsg.TaskId {
	case task.GAME_COMMAND:
		event.Value = taskMsg.Body
	}
	// put copy of event in DB for later audit
	if taskMsg.TaskId.IsForAudit() {
		err := Put(event)
		logrus.Debugf("event DB ID: %s", event.ID.Hex())
		if err != nil {
			return err
		}
	}
	return nil
}

func SendTaskToSlave(hostId primitive.ObjectID, user User, taskMsg task.Message) (err error) {
	/*
		first, decide if we need to trigger actions
	*/
	// skip action for empty hostId
	if hostId == primitive.NilObjectID {
		logrus.Warn("SendTaskToSlave() HostId is empty, ignoring")
		return
	}

	// skip action for nonexistent hosts
	host, err := GetHostById(hostId)
	if err != nil {
		logrus.Info("SendTaskToSlave() host not found, ignoring")
		return
	}

	// skip action for disconnected hosts
	// TODO retry few times before giving up
	availableHosts := events.SSE.Channels()
	channelName := fmt.Sprintf("/v1/events/%s", host.Token)
	if !lib.StringInSlice(channelName, availableHosts) {
		logrus.Warn("SendTaskToSlave() host dc, ignoring")
		return
	}

	// prepare message for sending to slave as a task
	channel, _ := events.SSE.GetChannel(channelName)
	taskMsgStr, err := taskMsg.Serialize()
	if err != nil {
		logrus.Error("SendTaskToSlave() serialization of task message failed")
		return
	}

	// log this if necessary
	err = SaveAction(taskMsg, user, hostId, "", "")
	if err != nil {
		logrus.Error("SendTaskToSlave() action save to DB failed")
		return
	}

	// send task to slave
	channel.SendMessage(sse.NewMessage("", taskMsgStr, taskMsg.TaskId.StringValue()))
	logrus.Infof("SendTaskToSlave() Task %v sent to host %v", taskMsg.TaskId.String(), host.ID.Hex())
	return
}

/*

func SendMongoEvent(hostId primitive.ObjectID, taskMsg task.Message) (err error) {
	event := &Event{
		HostId:      hostId,
		TaskMessage: taskMsg,
	}
	err = Put(event)
	if err != nil {
		return err
	}

	return
}

func ListenEventChangesFromMongo() {
	pipeline := bson.D{
		{
			Key: "$match",
			Value: bson.D{
				{Key: "operationType", Value: "insert"},
			},
		},
	}

	changeStream, err := DB.Collection(EventCollection).Watch(context.TODO(), mongo.Pipeline{pipeline})
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
*/
