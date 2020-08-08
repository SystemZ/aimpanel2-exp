package tasks

import (
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/lib/task"
	"gitlab.com/systemz/aimpanel2/slave/config"
	"gitlab.com/systemz/aimpanel2/slave/model"
	"time"
)

// tasks below will be eventually finished by wrapper
func WrapperTaskHandler(taskMsg task.Message) {
	switch taskMsg.TaskId {
	case task.GAME_COMMAND:
		GsCmd(taskMsg.GameServerID, taskMsg.Body)
	case task.GAME_STOP_SIGTERM:
		GsStop(taskMsg.GameServerID)
	case task.GAME_STOP_SIGKILL:
		GsKill(taskMsg.GameServerID)
	case task.GAME_RESTART:
		GsRestart(taskMsg)
	case task.GAME_METRICS_FREQUENCY:
		model.SendTask(config.REDIS_PUB_SUB_WRAPPER_CH, taskMsg)
	case task.GAME_SHUTDOWN:
		GsShutdown(taskMsg)
	}
}

func GsStartGame(taskMsg task.Message) {
	game, err := model.GetGsGame(taskMsg.GameServerID)
	if err != nil {
		logrus.Printf("Something went wrong when sending msg: %v", err)
		return
	}
	wrapperTask := task.Message{
		TaskId:       task.GAME_START,
		GameServerID: taskMsg.GameServerID,
		Game:         &game,
	}

	model.SendTask(config.REDIS_PUB_SUB_WRAPPER_CH, wrapperTask)
}

func GsCmd(gsId string, cmdStr string) {
	taskMsg := task.Message{
		TaskId:       task.GAME_COMMAND,
		GameServerID: gsId,
		Body:         cmdStr,
	}
	model.SendTask(config.REDIS_PUB_SUB_WRAPPER_CH, taskMsg)
}

func GsStop(gsId string) {
	taskMsg := task.Message{
		// FIXME other task IDs for user CLI actions
		TaskId:       task.GAME_STOP_SIGTERM,
		GameServerID: gsId,
	}
	model.SendTask(config.REDIS_PUB_SUB_WRAPPER_CH, taskMsg)
}

func GsKill(gsId string) {
	taskMsg := task.Message{
		// FIXME other task IDs for user CLI actions
		TaskId:       task.GAME_STOP_SIGKILL,
		GameServerID: gsId,
	}
	model.SendTask(config.REDIS_PUB_SUB_WRAPPER_CH, taskMsg)
}

func GsRestart(taskMsg task.Message) {
	model.SetGsRestart(taskMsg.GameServerID, 0)
	model.SetGsGame(taskMsg.GameServerID, taskMsg.Game)

	GsStop(taskMsg.GameServerID)
	model.SetGsRestart(taskMsg.GameServerID, 1)

	go func() {
		<-time.After(time.Duration(taskMsg.StopTimeout) * time.Second)

		val, err := model.GetGsRestart(taskMsg.GameServerID)
		if err != nil {
			return
		}

		if val == 1 {
			model.SetGsRestart(taskMsg.GameServerID, -1)
		}
	}()
}

func GsShutdown(taskMsg task.Message) {
	GsCmd(taskMsg.GameServerID, taskMsg.Game.StopCommand)

	go func() {
		<-time.After(time.Duration(taskMsg.Game.StopTimeout) * time.Second)

		running, err := model.GetGsRunning(taskMsg.GameServerID)
		if err != nil {
			return
		}

		if running == 1 {
			logrus.Infof("GS %s running - sending sigterm", taskMsg.GameServerID)
			GsStop(taskMsg.GameServerID)
		}
	}()

	go func() {
		<-time.After(time.Duration(taskMsg.Game.StopHardTimeout) * time.Second)

		running, err := model.GetGsRunning(taskMsg.GameServerID)
		if err != nil {
			return
		}

		if running == 1 {
			logrus.Infof("GS %s still running - sending sigkill", taskMsg.GameServerID)

			GsKill(taskMsg.GameServerID)
		}
	}()
}
