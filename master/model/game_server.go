package model

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
)

type GameServer struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty" example:"1238206236281802752"`

	// User assigned name
	Name string `bson:"name" json:"name" example:"Ultra MC Server"`

	// Host ID
	//
	// required: true
	HostId primitive.ObjectID `bson:"host_id" json:"host_id" example:"100112233-4455-6677-8899-aabbccddeeff"`

	// State
	// 0 off, 1 running
	State uint `bson:"state" json:"state" example:"0"`

	// Game ID
	GameId uint `bson:"game_id" json:"game_id" example:"1"`

	// Game Version
	GameVersion string `bson:"game_version" json:"game_version" example:"1.14.2"`

	// Game
	GameJson string `bson:"game_json" json:"game_json"`

	// Metric Frequency
	MetricFrequency int `bson:"metric_frequency" json:"metric_frequency" example:"30"`

	MetricMaxS int `bson:"metric_max_s" json:"metric_max_s" example:"30"`

	// Stop Timeout
	StopTimeout int `bson:"stop_timeout" json:"stop_timeout" example:"30"`

	// Custom launch cmd
	CustomCmdStart string `bson:"custom_cmd_start" json:"custom_cmd_start" example:"java -jar mc.jar"`

	// Ports to forward from container to host
	Ports *[]GamePort `bson:"ports" json:"ports"`
}

type GamePort struct {
	Protocol      string `json:"protocol,omitempty" bson:"protocol"`
	Host          string `json:"host,omitempty" bson:"host"`
	PortHost      int    `json:"port_host,omitempty" bson:"port_host"`
	PortContainer int    `json:"port_container,omitempty" bson:"port_container"`
}

func (g GameServer) SerializePorts() (res string) {
	var resPorts []string
	if g.Ports == nil {
		return
	}
	for _, port := range *g.Ports {
		resPorts = append(resPorts, fmt.Sprintf("%v:%v->%v/%v", port.Host, port.PortHost, port.PortContainer, port.Protocol))
	}
	res = strings.Join(resPorts, " ")
	return
}

func (g *GameServer) GetCollectionName() string {
	return GameServerCollection
}

func (g *GameServer) GetID() primitive.ObjectID {
	return g.ID
}

func (g *GameServer) SetID(id primitive.ObjectID) {
	g.ID = id
}

func GetGameServers() ([]GameServer, error) {
	var gameServers []GameServer

	cur, err := DB.Collection(GameServerCollection).Find(context.TODO(),
		bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var gameServer GameServer
		if err := cur.Decode(&gameServer); err != nil {
			return nil, err
		}
		gameServers = append(gameServers, gameServer)
	}

	return gameServers, nil
}

func GetGameServerById(id primitive.ObjectID) (*GameServer, error) {
	var gs GameServer
	err := DB.Collection(GameServerCollection).FindOne(context.TODO(), bson.D{{Key: "_id", Value: id}}).Decode(&gs)
	if err != nil {
		return nil, err
	}

	return &gs, nil
}

func GetGameServerByGsIdAndHostId(gsId primitive.ObjectID, hostId primitive.ObjectID) (*GameServer, error) {
	var gs GameServer

	err := DB.Collection(GameServerCollection).FindOne(context.TODO(), bson.D{
		{Key: "_id", Value: gsId},
		{Key: "host_id", Value: hostId},
	}).Decode(&gs)
	if err != nil {
		return nil, err
	}

	return &gs, nil
}

func GetGameServersByHostId(hostId primitive.ObjectID) (*[]GameServer, error) {
	var gameServers []GameServer

	cur, err := DB.Collection(GameServerCollection).Find(context.TODO(),
		bson.D{{Key: "host_id", Value: hostId}})
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var gameServer GameServer
		if err := cur.Decode(&gameServer); err != nil {
			return nil, err
		}
		gameServers = append(gameServers, gameServer)
	}

	return &gameServers, nil
}

//FIXME
func GetUserGameServers(user User) (*[]GameServer, error) {
	hosts, err := GetHostsByUser(user)
	if err != nil {
		return nil, err
	}

	var hostsId = make([]bson.D, 0)
	for _, host := range hosts {
		hostsId = append(hostsId, bson.D{{Key: "host_id", Value: host.ID}})
	}

	var gameServers = make([]GameServer, 0)
	if len(hostsId) == 0 {
		return &gameServers, nil
	}

	cur, err := DB.Collection(GameServerCollection).Find(context.TODO(),
		bson.D{{Key: "$or", Value: hostsId}})
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var gameServer GameServer
		if err := cur.Decode(&gameServer); err != nil {
			return nil, err
		}
		gameServers = append(gameServers, gameServer)
	}

	return &gameServers, nil
}
