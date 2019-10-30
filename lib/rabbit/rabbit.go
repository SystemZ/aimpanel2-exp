package rabbit

import (
	"github.com/gofrs/uuid"
	"gitlab.com/systemz/aimpanel2/lib/game"
	"gitlab.com/systemz/aimpanel2/master/model"
)

const (
	//WRAPPER
	GAME_START = iota
	GAME_COMMAND
	GAME_STOP_SIGKILL
	GAME_STOP_SIGTERM

	//AGENT
	WRAPPER_START
	GAME_INSTALL

	SERVER_LOG
	WRAPPER_STARTED
	WRAPPER_EXITED

	WRAPPER_METRICS_FREQUENCY
	WRAPPER_METRICS

	AGENT_METRICS_FREQUENCY
	AGENT_METRICS
	AGENT_OS
	AGENT_REMOVE_GS

	AGENT_HEARTBEAT
	WRAPPER_HEARTBEAT
)

type QueueMsg struct {
	// task id
	TaskId int `json:"task_id,omitempty"`

	// task specific arguments
	//Game             string               `json:"game,omitempty"`
	Game         game.Game       `json:"game,omitempty"`
	GameServerID uuid.UUID       `json:"game_server_id,omitempty"`
	Body         string          `json:"body,omitempty"`
	GameFile     *model.GameFile `json:"game_file,omitempty"`
	//GameCommands     *[]model.GameCommand `json:"game_commands,omitempty"`
	//GameStartCommand *model.GameCommand   `json:"game_start_command,omitempty"`

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
	RamTotal        int `json:"ram_total,omitempty"`
	DiskFree        int `json:"disk_free,omitempty"`
	DiskTotal       int `json:"disk_total,omitempty"`
	DiskUsed        int `json:"disk_used,omitempty"`

	AgentToken string `json:"agent_token,omitempty"`
	User       int    `json:"user,omitempty"`
	System     int    `json:"system,omitempty"`
	Idle       int    `json:"idle,omitempty"`
	Nice       int    `json:"nice,omitempty"`
	Iowait     int    `json:"iowait,omitempty"`
	Irq        int    `json:"irq,omitempty"`
	Softirq    int    `json:"softirq,omitempty"`
	Steal      int    `json:"steal,omitempty"`
	Guest      int    `json:"guest,omitempty"`
	GuestNice  int    `json:"guest_nice,omitempty"`

	OS              string `json:"os,omitempty"`
	Platform        string `json:"platform,omitempty"`
	PlatformFamily  string `json:"platform_family,omitempty"`
	PlatformVersion string `json:"platform_version,omitempty"`
	KernelVersion   string `json:"kernel_version,omitempty"`
	KernelArch      string `json:"kernel_arch,omitempty"`

	Timestamp int64 `json:"timestamp,omitempty"`
}

type LogMessage struct {
	Message string
}

type ExitMessage struct {
	Code    int
	Message string
}
