package rabbit

import (
	"github.com/gofrs/uuid"
	"gitlab.com/systemz/aimpanel2/master/model"
)

const (
	//WRAPPER
	GAME_START        = iota
	GAME_COMMAND      = iota
	GAME_STOP_SIGKILL = iota
	GAME_STOP_SIGTERM = iota

	//AGENT
	WRAPPER_START = iota
	GAME_INSTALL  = iota

	SERVER_LOG      = iota
	WRAPPER_STARTED = iota
	WRAPPER_EXITED  = iota

	WRAPPER_METRICS_FREQUENCY = iota
	WRAPPER_METRICS = iota
)

type QueueMsg struct {
	// task id
	TaskId int `json:"task_id,omitempty"`

	// task specific arguments
	Game             string               `json:"game,omitempty"`
	GameServerID     uuid.UUID            `json:"game_server_id,omitempty"`
	Body             string               `json:"body,omitempty"`
	GameFile         *model.GameFile      `json:"game_file,omitempty"`
	GameCommands     *[]model.GameCommand `json:"game_commands,omitempty"`
	GameStartCommand *model.GameCommand   `json:"game_start_command,omitempty"`

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

	MetricFrequency int `json:"metric_frequency,omitempty"`
	CpuUsage        int `json:"cpu_usage,omitempty"`
	RamUsage        int `json:"ram_usage,omitempty"`
	RamFree         int `json:"ram_free,omitempty"`
	DiskFree        int `json:"disk_free,omitempty"`
}

type LogMessage struct {
	Message string
}

type ExitMessage struct {
	Code    int
	Message string
}
