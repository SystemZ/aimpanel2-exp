package game_server

//swagger:parameters GameServer stop
type StopGameServerRequest struct {
	// required: true
	Type uint `json:"type"`
}
