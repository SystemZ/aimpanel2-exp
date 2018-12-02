package lib

const (
	START        = 0
	COMMAND      = 1
	STOP_SIGKILL = 2
	STOP_SIGTERM = 3
	DOWNLOAD     = 4
	OS_COMMAND   = 5
)

type RpcMessage struct {
	Type int
	Body string
}

type LogMessage struct {
	Message string
}

type ExitMessage struct {
	Code    int
	Message string
}

type Game struct {
	Name        string
	Command     string
	DownloadUrl string
	Path        string
}
