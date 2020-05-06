package model

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	STDOUT = iota
	STDERR = iota
)

type GameServerLog struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty" example:"1238206236281802752"`

	GameServerId primitive.ObjectID `json:"game_server_id"`

	Type uint `json:"type"`

	Log string `json:"log"`
}

func (g *GameServerLog) GetCollectionName() string {
	return "game_servers_log"
}

func GetLogsByGameServer(gsId string, limit int) *[]GameServerLog {
	var logs []GameServerLog

	err := GetSLimit(&logs, limit, map[string]string{
		"game_server_id": gsId,
	})
	if err != nil {
		return nil
	}

	return &logs
}
