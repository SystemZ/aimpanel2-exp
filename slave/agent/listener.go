package agent

import (
	"bytes"
	"github.com/inconshreveable/go-update"
	"github.com/r3labs/sse"
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/lib"
	"gitlab.com/systemz/aimpanel2/lib/task"
	"gitlab.com/systemz/aimpanel2/slave/config"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

func agent() {
	client := sse.NewClient(config.API_URL + "/v1/events/" + config.HOST_TOKEN)
	client.Headers = map[string]string{
		"Authorization": "Bearer " + config.API_TOKEN,
	}
	err := client.SubscribeRaw(func(msg *sse.Event) {
		logrus.Info(msg.ID)
		logrus.Info(string(msg.Data))
		logrus.Info(string(msg.Event))

		taskMsg := task.Message{}
		err := taskMsg.Deserialize(string(msg.Data))
		if err != nil {
			logrus.Error(err)
		}
		taskId, _ := strconv.Atoi(string(msg.Event))

		switch taskId {
		case task.WRAPPER_START:
			logrus.Info("START_WRAPPER")
			//TODO: move gsID to env variable
			cmd := exec.Command("slave", "wrapper", taskMsg.GameServerID)

			cmd.Env = os.Environ()
			cmd.Env = append(cmd.Env, "HOST_TOKEN="+config.HOST_TOKEN)
			cmd.Env = append(cmd.Env, "API_TOKEN="+config.API_TOKEN)

			//TODO: FOR TESTING ONLY
			var stdBuffer bytes.Buffer
			mw := io.MultiWriter(os.Stdout, &stdBuffer)
			cmd.Stdout = mw
			cmd.Stderr = mw

			if err := cmd.Start(); err != nil {
				logrus.Error(err)
			}

			cmd.Process.Release()
		case task.GAME_INSTALL:
			logrus.Info("INSTALL_GAME_SERVER")

			gsPath := filepath.Clean(config.GS_DIR) + "/" + taskMsg.GameServerID
			if _, err := os.Stat(gsPath); os.IsNotExist(err) {
				//TODO: Set correct perms
				_ = os.Mkdir(gsPath, 0777)
			}

			err = taskMsg.Game.Install(filepath.Clean(config.STORAGE_DIR), gsPath)
			if err != nil {
				logrus.Error(err)
			}

			logrus.Info("Installation finished")
		case task.AGENT_METRICS_FREQUENCY:
			logrus.Info("AGENT_METRICS_FREQUENCY")

			metricsFrequency = taskMsg.MetricFrequency

			go metrics()
		case task.AGENT_REMOVE_GS:
			logrus.Info("AGENT_REMOVE_GS")

			gsId := taskMsg.GameServerID
			gsPath := filepath.Clean(config.GS_DIR) + "/" + gsId
			gsTrashPath := filepath.Clean(config.TRASH_DIR) + "/" + gsId

			err := os.Rename(gsPath, gsTrashPath)
			if err != nil {
				logrus.Error(err)
			}
		case task.SLAVE_UPDATE:
			if config.GIT_COMMIT == taskMsg.Commit {
				return
			}

			resp, err := http.Get(taskMsg.Url)
			if err != nil {
				logrus.Error(err)
			}
			defer resp.Body.Close()

			err = update.Apply(resp.Body, update.Options{
				TargetPath:  "",
				TargetMode:  0,
				Checksum:    nil,
				PublicKey:   nil,
				Signature:   nil,
				Verifier:    nil,
				Hash:        0,
				Patcher:     nil,
				OldSavePath: "",
			})
			if err != nil {
				logrus.Error(err)
			}
			os.Exit(0)
		}
	})
	if err != nil {
		lib.FailOnError(err, "Failed to subscribe a channel")
	}
}
