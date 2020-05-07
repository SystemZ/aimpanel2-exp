package model

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	STDOUT = iota
	STDERR = iota
)

type GameServerLog struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty" example:"1238206236281802752"`

	GameServerId primitive.ObjectID `bson:"game_server_id" json:"game_server_id"`

	Type uint `bson:"type" json:"type"`

	Log string `bson:"log" json:"log"`
}

func (g *GameServerLog) GetCollectionName() string {
	return gameServerLogCollection
}

func (g *GameServerLog) GetID() primitive.ObjectID {
	return g.ID
}

func (g *GameServerLog) SetID(id primitive.ObjectID) {
	g.ID = id
}

func GetLogsByGameServerId(gsId primitive.ObjectID, limit int64) (*[]GameServerLog, error) {
	var logs []GameServerLog

	opts := options.Find()
	opts.SetLimit(limit)

	cur, err := DB.Collection(gameServerLogCollection).Find(context.TODO(),
		bson.D{{"game_server_id", gsId}}, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var log GameServerLog
		if err := cur.Decode(&log); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	return &logs, nil
}
