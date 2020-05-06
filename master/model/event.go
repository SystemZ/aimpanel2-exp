package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Event struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty" example:"1238206236281802752"`

	EventId int `json:"event_id"`

	HostId string `json:"host_id"`

	UserId string `json:"user_id"`

	GameServerId string `json:"game_server_id"`
}
