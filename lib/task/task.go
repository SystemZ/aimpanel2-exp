package task

import (
	"encoding/json"
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

type Message struct {
	// task id
	TaskId       int             `json:"task_id,omitempty"`
	Game         game.Game       `json:"game,omitempty"`
	GameServerID uuid.UUID       `json:"game_server_id,omitempty"`
	Body         string          `json:"body,omitempty"`
	GameFile     *model.GameFile `json:"game_file,omitempty"`
	AgentToken   string          `json:"agent_token,omitempty"`

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
