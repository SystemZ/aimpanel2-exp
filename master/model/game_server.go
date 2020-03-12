package model

import (
	"github.com/jinzhu/gorm"
	"time"
)

type GameServer struct {
	// User assigned name
	Name string `json:"name" example:"Ultra MC Server"`

	// Host ID
	//
	// required: true
	HostId string `json:"host_id" example:"100112233-4455-6677-8899-aabbccddeeff"`

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

func GetGameServer(gsId string) *GameServer {
	var gs GameServer

	err := GetOneS(&gs, map[string]interface{}{
		"doc_type": "game_server",
		"id":       gsId,
	})
	if err != nil {
		return nil
	}

	return &gs
}

func GetGameServerByGsIdAndHostId(serverId string, hostId string) *GameServer {
	var gs GameServer

	err := GetOneS(&gs, map[string]interface{}{
		"doc_type": "game_server",
		"id":       serverId,
		"hostId":   hostId,
	})
	if err != nil {
		return nil
	}

	return &gs
}

func GetGameServersByHostId(hostId string) *[]GameServer {
	var gs []GameServer

	err := GetOneS(&gs, map[string]interface{}{
		"doc_type": "game_server",
		"hostId":   hostId,
	})
	if err != nil {
		return nil
	}

	return &gs
}

//FIXME
func GetUserGameServers(db *gorm.DB, userId string) *[]GameServer {
	//var gameServers []GameServer
	//
	//if db.Table("game_servers").Select("game_servers.*").Joins(
	//	"LEFT JOIN hosts ON game_servers.host_id = hosts.id").Where(
	//	"hosts.user_id = ?", userId).Find(&gameServers).RecordNotFound() {
	//	return nil
	//}
	//
	//return &gameServers
	return nil
}
