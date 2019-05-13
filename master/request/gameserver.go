package request

type GameServerStopReq struct {
	Type uint `json:"type"`
}

type GameServerSendCommandReq struct {
	Command string `json:"command"`
}
