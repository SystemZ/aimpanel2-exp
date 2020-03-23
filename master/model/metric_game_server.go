package model

type MetricGameServer struct {
	Base

	GameServerId string `json:"game_server_id"`

	RamUsage int `json:"ram_usage"`

	CpuUsage int `json:"cpu_usage"`
}
