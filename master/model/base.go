package model

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	userCollection             = "user"
	groupCollection            = "group"
	groupUserCollection        = "group_user"
	permissionCollection       = "permission"
	hostJobCollection          = "host_job"
	hostCollection             = "host"
	metricHostCollection       = "metric_host"
	metricGameServerCollection = "metric_game_server"
	gameServerCollection       = "game_server"
	gameServerLogCollection    = "game_server_log"
	gameFileCollection         = "game_file"
)

type Document interface {
	GetID() primitive.ObjectID
	SetID(id primitive.ObjectID)
	GetCollectionName() string
}

func Put(d Document) error {
	res, err := DB.Collection(d.GetCollectionName()).InsertOne(context.TODO(), d)
	if err != nil {
		return err
	}

	d.SetID(res.InsertedID.(primitive.ObjectID))

	return nil
}

func Update(d Document) error {
	_, err := DB.Collection(d.GetCollectionName()).UpdateOne(
		context.TODO(),
		bson.D{{"_id", d.GetID()}},
		bson.D{{"$set", d}})
	if err != nil {
		return err
	}

	return nil
}

func Delete(d Document) error {
	_, err := DB.Collection(d.GetCollectionName()).DeleteOne(
		context.TODO(),
		bson.D{{"_id", d.GetID()}})
	if err != nil {
		return err
	}

	return nil
}
