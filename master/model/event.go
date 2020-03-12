package model

type Event struct {
	Base

	EventId int `json:"event_id"`

	HostId string `json:"host_id"`

	UserId string `json:"user_id"`

	GameServerId string `json:"game_server_id"`
}
