package task

import (
	"encoding/json"
	"gitlab.com/systemz/aimpanel2/lib/game"
	"gitlab.com/systemz/aimpanel2/master/model"
)

const (
	//WRAPPER
	GAME_START = iota + 1
	GAME_COMMAND
	GAME_STOP_SIGKILL
	GAME_STOP_SIGTERM

	SLAVE_UPDATE

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

type Message struct {
	// task id
	TaskId       int             `json:"task_id,omitempty"`
	Game         game.Game       `json:"game,omitempty"`
	GameServerID string          `json:"game_server_id,omitempty"`
	Body         string          `json:"body,omitempty"`
	GameFile     *model.GameFile `json:"game_file,omitempty"`

	// task progress
	Stdout string `json:"stdout,omitempty"`
	Stderr string `json:"stderr,omitempty"`

	MetricFrequency int `json:"metric_frequency,omitempty"`
	CpuUsage        int `json:"cpu_usage,omitempty"`
	RamUsage        int `json:"ram_usage,omitempty"`
	RamFree         int `json:"ram_free,omitempty"`
	RamTotal        int `json:"ram_total,omitempty"`
	DiskFree        int `json:"disk_free,omitempty"`
	DiskTotal       int `json:"disk_total,omitempty"`
	DiskUsed        int `json:"disk_used,omitempty"`

	User      int `json:"user,omitempty"`
	System    int `json:"system,omitempty"`
	Idle      int `json:"idle,omitempty"`
	Nice      int `json:"nice,omitempty"`
	Iowait    int `json:"iowait,omitempty"`
	Irq       int `json:"irq,omitempty"`
	Softirq   int `json:"softirq,omitempty"`
	Steal     int `json:"steal,omitempty"`
	Guest     int `json:"guest,omitempty"`
	GuestNice int `json:"guest_nice,omitempty"`

	OS              string `json:"os,omitempty"`
	Platform        string `json:"platform,omitempty"`
	PlatformFamily  string `json:"platform_family,omitempty"`
	PlatformVersion string `json:"platform_version,omitempty"`
	KernelVersion   string `json:"kernel_version,omitempty"`
	KernelArch      string `json:"kernel_arch,omitempty"`

	Commit string `json:"commit,omitempty"`
	Url    string `json:"url,omitempty"`

	Timestamp int64 `json:"timestamp,omitempty"`
}

func (m *Message) Serialize() (string, error) {
	data, err := json.Marshal(&m)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (m *Message) Deserialize(data string) error {
	err := json.Unmarshal([]byte(data), &m)
	if err != nil {
		return err
	}
	return nil
}
