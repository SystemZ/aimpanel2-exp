package tasks

import (
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/lib/task"
	"gitlab.com/systemz/aimpanel2/slave/config"
	"gitlab.com/systemz/aimpanel2/slave/model"
	"os"
	"syscall"
)

// tasks below will be eventually finished by agent
func StartWrapperExecRaw(taskMsg task.Message) {
	cred := &syscall.Credential{
		Uid:         uint32(syscall.Getuid()),
		Gid:         uint32(syscall.Getgid()),
		Groups:      []uint32{},
		NoSetGroups: true,
	}

	sysproc := &syscall.SysProcAttr{
		Credential: cred,
		Noctty:     true,
		Setpgid:    true,
	}

	attr := os.ProcAttr{
		Dir: ".",
		Env: os.Environ(),
		Files: []*os.File{
			os.Stdin,
			os.Stdout,
			os.Stderr,
		},
		Sys: sysproc,
	}

	// TODO test this for very long running GS
	attr.Env = append(attr.Env, "HOST_TOKEN="+config.HOST_TOKEN)
	attr.Env = append(attr.Env, "API_TOKEN="+config.API_TOKEN)

	//TODO: move gsID to env variable
	//FIXME use variables, don't hardcode paths
	process, err := os.StartProcess("/usr/local/bin/slave", []string{"/usr/local/bin/slave", "wrapper", taskMsg.GameServerID}, &attr)
	if err != nil {
		logrus.Error(err)
	}

	model.SetGsRunning(taskMsg.GameServerID, 1)

	go func() {
		_, err := process.Wait()
		if err != nil {
			logrus.Error(err)
		}

		taskMsg := task.Message{
			TaskId:       task.GAME_SHUTDOWN,
			GameServerID: taskMsg.GameServerID,
		}
		model.SendTask(config.REDIS_PUB_SUB_AGENT_CH, taskMsg)
		model.SetGsRunning(taskMsg.GameServerID, 0)
	}()
}
