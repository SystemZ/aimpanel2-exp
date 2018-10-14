package lib

const (
	START        = 0
	COMMAND      = 1
	STOP_SIGKILL = 2
	STOP_SIGTERM = 3
)

type RpcMessage struct {
	Type int
	Body string
}

type LogMessage struct {
	Message string
}
