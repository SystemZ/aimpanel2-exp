package request

type HostCreateRequest struct {
	// User assigned name
	Name string `json:"name" example:"My Great Linux server"`

	// User assigned ip
	Ip string `json:"ip" example:"192.51.100.128"`
}

type GameServerStopRequest struct {
	Type uint `json:"type" example:"1"`
}

type GameServerSendCommandRequest struct {
	Command string `json:"command" example:"say Hello!"`
}
