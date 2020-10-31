package response

import (
	"github.com/gofrs/uuid"
	"gitlab.com/systemz/aimpanel2/master/model"
)

type Token struct {
	Token string `json:"token"`
}

type JsonSuccess struct {
	Message string `json:"message" example:""`
}

//This need to be here for swagger
type JsonError struct {
	ErrorCode int    `json:"error_code" example:"1"`
	Message   string `json:"message" example:""`
}

type Game struct {
	//ID of the game
	Id uint `json:"id" example:"1"`
	//Game name
	Name string `json:"name" example:"Spigot"`
	//Supported versions
	Versions []string `json:"versions" example:"1.14.1,1.12.2"`
}

type UserProfile struct {
	//ID of the user
	ID uuid.UUID `json:"id" example:"100112233-4455-6677-8899-aabbccddeeff"`

	//User assigned username
	Username string `json:"username" example:"john.doe"`

	//User assigned email
	Email string `json:"email" example:"john.doe@example.com"`
}

type UserProfileResponse struct {
	User UserProfile `json:"user"`
}

type Host struct {
	//Host details
	Host model.Host `json:"host"`
}

type HostList struct {
	//List of hosts
	Hosts []model.Host `json:"hosts"`
}

type HostMetrics struct {
	//Metric info
	Metrics []model.TimeseriesOutput `json:"metrics"`
}

type HostJobList struct {
	//List of jobs
	Jobs []model.HostJob `json:"jobs"`
}

type GameServerList struct {
	//List of game servers
	GameServers []model.GameServer `json:"game_servers"`
}

type GameServer struct {
	//Game server details
	GameServer model.GameServer `json:"game_server"`
}

type ID struct {
	ID string `json:"id"`
}

type BackupList struct {
	Backups []string `json:"backups"`
}
