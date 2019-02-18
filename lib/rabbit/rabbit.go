package rabbit

const (
	//WRAPPER
	GAME_START        = iota
	GAME_COMMAND      = iota
	GAME_STOP_SIGKILL = iota
	GAME_STOP_SIGTERM = iota

	//AGENT
	WRAPPER_START = iota
	GAME_INSTALL  = iota
)

type QueueMsg struct {
	// task id
	TaskId int

	// task specific arguments
	Game         string `json:"game,omitempty"`
	GameServerID string `json:"game_server_id,omitempty"`
	Body         string `json:"body,omitempty"`

	// task started
	TaskStarted bool `json:"task_started,omitempty"`
	// task progress
	Stdout string `json:"stdout,omitempty"`
	Stderr string `json:"stderr,omitempty"`
	// task end
	TaskResult string `json:"task_result,omitempty"`
	TaskEnd    bool   `json:"task_end,omitempty"`
	TaskOk     bool   `json:"task_ok,omitempty"`
	TaskError  bool   `json:"task_error,omitempty"`
}

type LogMessage struct {
	Message string
}

type ExitMessage struct {
	Code    int
	Message string
}
