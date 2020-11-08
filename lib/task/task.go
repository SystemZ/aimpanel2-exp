package task

import (
	"encoding/json"
	"gitlab.com/systemz/aimpanel2/lib/filemanager"
	"gitlab.com/systemz/aimpanel2/lib/game"
	"strconv"
)

type Id int

//go:generate stringer -type=Id
const (
	//WRAPPER
	GAME_START Id = iota + 1
	GAME_COMMAND
	GAME_STOP_SIGKILL
	GAME_STOP_SIGTERM
	GAME_RESTART
	GAME_SERVER_LOG
	GAME_STARTED
	GAME_SHUTDOWN
	GAME_METRICS_FREQUENCY
	GAME_METRICS

	AGENT_STARTED
	AGENT_SHUTDOWN
	AGENT_UPDATE
	AGENT_OS
	AGENT_METRICS
	AGENT_METRICS_FREQUENCY
	AGENT_GET_JOBS
	AGENT_GET_UPDATE

	AGENT_REMOVE_GS
	AGENT_BACKUP_GS
	AGENT_BACKUP_RESTORE_GS
	AGENT_BACKUP_LIST_GS
	AGENT_START_GS
	AGENT_INSTALL_GS
	AGENT_FILE_LIST_GS
	AGENT_CLEAN_REINSTALL_GS

	// special case, eg. just to keep up SSE session
	PING

	GS_CMD_START_CHANGE
	HOST_NAME_CHANGE
	GS_NAME_CHANGE
	GS_GAME_CHANGE
	HOST_HWID_CHANGE
)

func (i Id) StringValue() string {
	return strconv.Itoa(int(i))
}

func (i Id) IsForAudit() bool {
	switch i {
	case
		AGENT_FILE_LIST_GS,
		AGENT_GET_JOBS,
		GAME_METRICS_FREQUENCY,
		AGENT_METRICS_FREQUENCY:
		return false
	default:
		return true
	}
}

type Messages []Message

type Message struct {
	// task id
	TaskId             Id         `json:"task_id,omitempty"`
	Game               *game.Game `json:"game,omitempty"`
	GameServerID       string     `json:"game_server_id,omitempty"`
	HostID             string     `json:"host_id,omitempty"`
	GameCustomCmdStart string     `json:"game_custom_cmd_start"`
	Body               string     `json:"body,omitempty"`
	StopTimeout        int        `json:"stop_timeout,omitempty"`

	// task progress
	Stdout string `json:"stdout,omitempty"`
	Stderr string `json:"stderr,omitempty"`

	MetricFrequency int `json:"metric_frequency,omitempty"`
	CpuUsage        int `json:"cpu_usage,omitempty"`
	RamUsage        int `json:"ram_usage,omitempty"`
	RamFree         int `json:"ram_free,omitempty"`
	RamCache        int `json:"ram_cache,omitempty"`
	RamBuffers      int `json:"ram_buffers,omitempty"`
	RamAvailable    int `json:"ram_available,omitempty"`
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

	Files *filemanager.Node `json:"files,omitempty"`
	Jobs  *[]Job            `json:"jobs,omitempty"`
	Ports *[]Port           `json:"ports,omitempty"`

	BackupFilename string   `json:"backup_filename,omitempty"`
	Backups        []string `json:"backups,omitempty"`

	Timestamp int64 `json:"timestamp,omitempty"`
}

type Job struct {
	Name           string  `json:"name,omitempty"`
	CronExpression string  `json:"cron_expression,omitempty"`
	TaskMessage    Message `json:"task_message,omitempty"`
}

type Port struct {
	Protocol      string `json:"protocol,omitempty"`
	Host          string `json:"host,omitempty"`
	PortHost      int    `json:"port_host,omitempty"`
	PortContainer int    `json:"port_container,omitempty"`
}

func (m *Message) Serialize() (string, error) {
	data, err := json.Marshal(&m)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (m *Messages) Serialize() (string, error) {
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
