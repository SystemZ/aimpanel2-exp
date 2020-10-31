package tasks

import (
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/lib/task"
)

func ProcessTask(taskMsg task.Message) {
	switch taskMsg.TaskId {
	case task.PING:
		// empty, don't take any action
		// we don't need to log it
	case task.GAME_COMMAND, task.GAME_STOP_SIGKILL,
		task.GAME_STOP_SIGTERM, task.GAME_RESTART, task.GAME_METRICS_FREQUENCY,
		task.GAME_SHUTDOWN:
		logrus.Infof("wrapper got task %v", taskMsg.TaskId.String())
		// executed by wrapper
		WrapperTaskHandler(taskMsg)
	case task.AGENT_START_GS, task.AGENT_INSTALL_GS,
		task.AGENT_BACKUP_GS, task.AGENT_UPDATE,
		task.AGENT_REMOVE_GS, task.AGENT_FILE_LIST_GS,
		task.AGENT_METRICS_FREQUENCY, task.AGENT_GET_UPDATE,
		task.AGENT_BACKUP_RESTORE_GS:
		logrus.Infof("agent got task %v", taskMsg.TaskId.String())
		// executed by agent
		AgentTaskHandler(taskMsg)
	default:
		logrus.Infof("Unknown task %v!", taskMsg.TaskId)
	}
}
