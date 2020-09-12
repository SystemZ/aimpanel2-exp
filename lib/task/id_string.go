// Code generated by "stringer -type=Id lib/task/task.go"; DO NOT EDIT.

package task

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[GAME_START-1]
	_ = x[GAME_COMMAND-2]
	_ = x[GAME_STOP_SIGKILL-3]
	_ = x[GAME_STOP_SIGTERM-4]
	_ = x[GAME_RESTART-5]
	_ = x[GAME_SERVER_LOG-6]
	_ = x[GAME_STARTED-7]
	_ = x[GAME_SHUTDOWN-8]
	_ = x[GAME_METRICS_FREQUENCY-9]
	_ = x[GAME_METRICS-10]
	_ = x[AGENT_STARTED-11]
	_ = x[AGENT_SHUTDOWN-12]
	_ = x[AGENT_UPDATE-13]
	_ = x[AGENT_OS-14]
	_ = x[AGENT_METRICS-15]
	_ = x[AGENT_METRICS_FREQUENCY-16]
	_ = x[AGENT_GET_JOBS-17]
	_ = x[AGENT_GET_UPDATE-18]
	_ = x[AGENT_REMOVE_GS-19]
	_ = x[AGENT_BACKUP_GS-20]
	_ = x[AGENT_START_GS-21]
	_ = x[AGENT_INSTALL_GS-22]
	_ = x[AGENT_FILE_LIST_GS-23]
	_ = x[PING-24]
	_ = x[GS_CMD_START_CHANGE-25]
}

const _Id_name = "GAME_STARTGAME_COMMANDGAME_STOP_SIGKILLGAME_STOP_SIGTERMGAME_RESTARTGAME_SERVER_LOGGAME_STARTEDGAME_SHUTDOWNGAME_METRICS_FREQUENCYGAME_METRICSAGENT_STARTEDAGENT_SHUTDOWNAGENT_UPDATEAGENT_OSAGENT_METRICSAGENT_METRICS_FREQUENCYAGENT_GET_JOBSAGENT_GET_UPDATEAGENT_REMOVE_GSAGENT_BACKUP_GSAGENT_START_GSAGENT_INSTALL_GSAGENT_FILE_LIST_GSPINGGS_CMD_START_CHANGE"

var _Id_index = [...]uint16{0, 10, 22, 39, 56, 68, 83, 95, 108, 130, 142, 155, 169, 181, 189, 202, 225, 239, 255, 270, 285, 299, 315, 333, 337, 356}

func (i Id) String() string {
	i -= 1
	if i < 0 || i >= Id(len(_Id_index)-1) {
		return "Id(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _Id_name[_Id_index[i]:_Id_index[i+1]]
}
