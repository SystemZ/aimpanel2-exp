package model

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type GameServer struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty" example:"1238206236281802752"`

	// User assigned name
	Name string `json:"name" example:"Ultra MC Server"`

	// Host ID
	//
	// required: true
	HostId primitive.ObjectID `json:"host_id" example:"100112233-4455-6677-8899-aabbccddeeff"`

	// State
	// 0 off, 1 running
	State uint `json:"state" example:"0"`

	// State Last Changed
	//FIXME default current timestamp
	StateLastChanged time.Time `json:"state_last_changed" example:"2019-09-29T03:16:27+02:00"`

	// Game ID
	GameId uint `json:"game_id" example:"1"`

	// Game Version
	GameVersion string `json:"game_version" example:"1.14.2"`

	// Game
	GameJson string `json:"game_json"`

	// Metric Frequency
	MetricFrequency int `json:"metric_frequency" example:"30"`

	// Stop Timeout
	StopTimeout int `json:"stop_timeout" example:"30"`
}

func (g *GameServer) GetCollectionName() string {
	return "game_servers"
}

func (g *GameServer) GetID() primitive.ObjectID {
	return g.ID
}

func GetGameServers() ([]GameServer, error) {
	var gameServers []GameServer

	cur, err := DB.Collection(gameServerCollection).Find(context.TODO(),
		bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var gameServer GameServer
		if err := cur.Decode(gameServer); err != nil {
			return nil, err
		}
		gameServers = append(gameServers, gameServer)
	}

	return gameServers, nil
}

func GetGameServerById(id primitive.ObjectID) (*GameServer, error) {
	var gs GameServer

	err := DB.Collection(gameServerCollection).FindOne(context.TODO(), bson.D{{"_id", id}}).Decode(&gs)
	if err != nil {
		return nil, err
	}

	return &gs, nil
}

func GetGameServerByGsIdAndHostId(gsId primitive.ObjectID, hostId primitive.ObjectID) (*GameServer, error) {
	var gs GameServer

	err := DB.Collection(gameServerCollection).FindOne(context.TODO(), bson.D{
		{"_id", gsId},
		{"host_id", hostId},
	}).Decode(&gs)
	if err != nil {
		return nil, err
	}

	return &gs, nil
}

func GetGameServersByHostId(hostId primitive.ObjectID) (*[]GameServer, error) {
	var gameServers []GameServer

	cur, err := DB.Collection(gameServerCollection).Find(context.TODO(),
		bson.D{{"host_id", hostId}})
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var gameServer GameServer
		if err := cur.Decode(gameServer); err != nil {
			return nil, err
		}
		gameServers = append(gameServers, gameServer)
	}

	return &gameServers, nil
}

//FIXME
func GetUserGameServers(userId primitive.ObjectID) (*[]GameServer, error) {
	hosts, err := GetHostsByUserId(userId)
	if err != nil {
		return nil, err
	}

	var hostsId []bson.D
	for _, host := range hosts {
		hostsId = append(hostsId, bson.D{{"host_id", host.ID}})
	}

	var gameServers []GameServer

	cur, err := DB.Collection(gameServerCollection).Find(context.TODO(),
		bson.D{{"$or", hostsId}})
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var gameServer GameServer
		if err := cur.Decode(gameServer); err != nil {
			return nil, err
		}
		gameServers = append(gameServers, gameServer)
	}

	return &gameServers, nil
}
