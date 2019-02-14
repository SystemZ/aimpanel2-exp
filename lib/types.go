package lib

const (
	//WRAPPER
	GAME_START        = iota
	GAME_COMMAND      = iota
	GAME_STOP_SIGKILL = iota
	GAME_STOP_SIGTERM = iota

	//AGENT
	WRAPPER_START = iota
	GAME_INSTALL  = iota
	GAME_RESTART  = iota
)

type RpcMessage struct {
	Type           int
	Body           string
	Game           string
	GameServerUUID string
}

type LogMessage struct {
	Message string
}

type ExitMessage struct {
	Code    int
	Message string
}
