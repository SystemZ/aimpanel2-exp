package model

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	UserCollection             = "user"
	PermissionCollection       = "user_permission"
	HostJobCollection          = "host_job"
	HostCollection             = "host"
	MetricHostCollection       = "metric_host"
	MetricGameServerCollection = "metric_game_server"
	GameServerCollection       = "game_server"
	GameServerLogCollection    = "game_server_log"
	GameFileCollection         = "game_file"
	MetricCollection           = "metrics"
	EventCollection            = "event"
	ConfigCollection           = "config"

	//LE
	CertDomainCollection = "cert_domains"
	CertCollection       = "certs"
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
		bson.D{{Key: "_id", Value: d.GetID()}},
		bson.D{{Key: "$set", Value: d}})
	if err != nil {
		return err
	}

	return nil
}

/*
func Upsert(d Document) error {
	upsert := true
	_, err := DB.Collection(d.GetCollectionName()).UpdateOne(
		context.TODO(),
		bson.D{{Key: "$set", Value: d}},
		&options.UpdateOptions{
			Upsert: &upsert,
		},
	)
	if err != nil {
		return err
	}

	return nil
}
*/

func Delete(d Document) error {
	_, err := DB.Collection(d.GetCollectionName()).DeleteOne(
		context.TODO(),
		bson.D{{Key: "_id", Value: d.GetID()}})
	if err != nil {
		return err
	}

	return nil
}
