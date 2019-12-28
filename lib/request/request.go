package request

type HostCreate struct {
	// User assigned name
	Name string `json:"name" example:"My Great Linux server"`

	// User assigned ip
	Ip string `json:"ip" example:"192.51.100.128"`
}

type GameServerCreate struct {
	//User assigned name
	Name string `json:"name" example:"Ultra MC Server"`

	//Selected game id
	GameId uint `json:"game_id" example:"1"`

	//Selected game version
	GameVersion string `json:"game_version" example:"1.14.2"`
}

type GameServerStop struct {
	Type uint `json:"type" example:"1"`
}

type GameServerSendCommand struct {
	Command string `json:"command" example:"say Hello!"`
}
