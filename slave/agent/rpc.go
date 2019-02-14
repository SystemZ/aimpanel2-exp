package agent

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"gitlab.com/systemz/aimpanel2/lib"
	"os"
	"os/exec"
	"strings"
)

func Rpc(channel *amqp.Channel, rpcQueue amqp.Queue) {
	msgs, err := channel.Consume(rpcQueue.Name, "", false, false, false, false, nil)
	lib.FailOnError(err, "Failed to register a consumer")

	for d := range msgs {
		logrus.Info("Got RPC call from RabbitMQ")

		var rpcMsg lib.RpcMessage
		err := json.Unmarshal(d.Body, &rpcMsg)
		if err != nil {
			logrus.Warn(err)
		}

		switch rpcMsg.Type {
		case lib.GAME_INSTALL:
			logrus.Info("INSTALL_GAME_SERVER")

			game := lib.GAMES[rpcMsg.Game]

			logrus.Info("Creating gs dir")

			gsPath := "/opt/aimpanel/gs/" + rpcMsg.GameServerUUID
			if _, err := os.Stat(gsPath); os.IsNotExist(err) {
				os.Mkdir(gsPath, 0777)
			}

			logrus.Info("Downloading install package")

			if _, err = os.Stat("/opt/aimpanel/storage/" + game.FileName); os.IsNotExist(err) {
				cmd := exec.Command("wget", game.DownloadUrl)
				cmd.Dir = "/opt/aimpanel/storage"

				if err := cmd.Run(); err != nil {
					logrus.Error(err)
				}

				cmd.Wait()
			}

			logrus.Info("Executing install commands")

			for _, c := range game.InstallCmds {
				var command []string
				for _, arg := range c {
					arg = strings.Replace(arg, "{uuid}", rpcMsg.GameServerUUID, -1)
					arg = strings.Replace(arg, "{fileName}", game.FileName, -1)

					command = append(command, arg)
				}

				logrus.Info("Executing")
				logrus.Info(command)

				cmd := exec.Command(command[0], command[1:]...)
				cmd.Dir = gsPath

				if err = cmd.Run(); err != nil {
					logrus.Error(err)
				}

				cmd.Wait()
			}

			logrus.Info("Installation finished")

			d.Ack(false)
		case lib.WRAPPER_START:
			logrus.Info("START_WRAPPER")
			cmd := exec.Command("slave", "wrapper", "test-test-test-test")
			if err := cmd.Start(); err != nil {
				logrus.Error(err)
			}
			cmd.Process.Release()

			d.Ack(false)
		}
	}
}
