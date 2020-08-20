package wrapper

import (
	"gitlab.com/systemz/aimpanel2/lib/task"
	"gitlab.com/systemz/aimpanel2/slave/config"
	"gitlab.com/systemz/aimpanel2/slave/model"
)

func (p *Process) LogStdout(msg string) {
	taskMsg := task.Message{
		TaskId:       task.GAME_SERVER_LOG,
		GameServerID: p.GameServerID,
		Stdout:       msg,
	}
	model.SendTask(config.REDIS_PUB_SUB_AGENT_CH, taskMsg)
}

func (p *Process) LogStderr(msg string) {
	taskMsg := task.Message{
		TaskId:       task.GAME_SERVER_LOG,
		GameServerID: p.GameServerID,
		Stderr:       msg,
	}
	model.SendTask(config.REDIS_PUB_SUB_AGENT_CH, taskMsg)
}
