package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type MetricGameServer struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty" example:"1238206236281802752"`

	GameServerId primitive.ObjectID `bson:"game_server_id" json:"game_server_id"`

	RamUsage int `bson:"ram_usage" json:"ram_usage"`

	CpuUsage int `bson:"cpu_usage" json:"cpu_usage"`
}

func (m *MetricGameServer) GetCollectionName() string {
	return metricGameServerCollection
}

func (m *MetricGameServer) GetID() primitive.ObjectID {
	return m.ID
}

func (m *MetricGameServer) SetID(id primitive.ObjectID) {
	m.ID = id
}
