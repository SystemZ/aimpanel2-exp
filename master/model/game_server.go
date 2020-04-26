package model

import (
	"time"
)

type GameServer struct {
	Base
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

func GetGameServers() []GameServer {
	var gs []GameServer

	err := GetOneS(&gs, map[string]interface{}{
		"doc_type": "game_server",
	})
	if err != nil {
		return nil
	}

	return gs
}

func GetGameServer(gsId string) *GameServer {
	var gs GameServer

	err := GetOneS(&gs, map[string]interface{}{
		"doc_type": "game_server",
		"_id":      gsId,
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
		"_id":      serverId,
		"host_id":  hostId,
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
		"host_id":  hostId,
	})
	if err != nil {
		return nil
	}

	return &gs
}

//FIXME
func GetUserGameServers(userId string) *[]GameServer {
	hosts := GetHostsByUserId(userId)
	var hostsId []interface{}
	for _, host := range hosts {
		hostsId = append(hostsId, map[string]interface{}{
			"host_id": host.ID,
		})
	}

	var gameServers []GameServer
	err := GetS(&gameServers, map[string]interface{}{
		"doc_type": "game_server",
		"$or":      hostsId,
	})
	if err != nil {
		return nil
	}

	return &gameServers
}
