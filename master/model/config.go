package model

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Config struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty" example:"1238206236281802752"`

	MigrateVersion int `bson:"migrate_version" json:"migrate_version"`
}

func (u *Config) GetCollectionName() string {
	return configCollection
}

func (u *Config) GetID() primitive.ObjectID {
	return u.ID
}

func (u *Config) SetID(id primitive.ObjectID) {
	u.ID = id
}

func GetConfig() (config Config, err error) {
	err = DB.Collection(configCollection).FindOne(context.TODO(), bson.D{}).Decode(&config)
	return
}

func IsConfigDocPresent() (exists bool, err error) {
	count, err := DB.Collection(configCollection).CountDocuments(context.TODO(), bson.D{})
	if err != nil {
		logrus.Error(err)
		return
	}

	if count == 1 {
		exists = true
	} else if count > 1 {
		exists = true
		err = errors.New("more than one config doc")
		return
	}
	return
}
