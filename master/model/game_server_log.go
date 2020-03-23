package model

const (
	STDOUT = iota
	STDERR = iota
)

type GameServerLog struct {
	Base

	GameServerId string `json:"game_server_id"`

	Type uint `json:"type"`

	Log string `json:"log"`
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
