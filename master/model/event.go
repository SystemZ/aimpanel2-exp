package model

import (
	"context"
	"gitlab.com/systemz/aimpanel2/lib/task"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type Event struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty" example:"1238206236281802752"`

	TaskMessage task.Message `bson:"task_message" json:"task_message"`
}

func EventChanges() {
	changeStream, err := DB.Collection(eventCollection).Watch(context.TODO(), mongo.Pipeline{})
	if err != nil {
		log.Fatal(err)
	}
	defer changeStream.Close(context.TODO())

	for changeStream.Next(context.TODO()) {
		var event bson.M
		if err := changeStream.Decode(&event); err != nil {
			log.Fatal(err)
		}

		log.Println(event)
	}

	if err := changeStream.Err(); err != nil {
		log.Fatal(err)
	}
}
