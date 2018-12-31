package lib

const (
	//WRAPPER
	START        = 0
	COMMAND      = 1
	STOP_SIGKILL = 2
	STOP_SIGTERM = 3

	//AGENT
	INSTALL_GAME_SERVER = 4
	START_WRAPPER       = 5
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
