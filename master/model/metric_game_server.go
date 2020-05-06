package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type MetricGameServer struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty" example:"1238206236281802752"`

	GameServerId primitive.ObjectID `json:"game_server_id"`

	RamUsage int `json:"ram_usage"`

	CpuUsage int `json:"cpu_usage"`
}

func (m *MetricGameServer) GetCollectionName() string {
	return "metrics_game_server"
}
