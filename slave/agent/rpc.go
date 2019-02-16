package agent

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"gitlab.com/systemz/aimpanel2/lib"
	"gitlab.com/systemz/aimpanel2/lib/rabbit"
	"gitlab.com/systemz/aimpanel2/slave/config"
	"os"
	"os/exec"
	"strings"
)

type rabbitTask struct {
	msg     amqp.Delivery
	msgBody rabbit.QueueMsg
	ch      *amqp.Channel
}

func rabbitListen(queue string) {
	conn, err := amqp.Dial("amqp://" + config.RABBITMQ_USERNAME + ":" + config.RABBITMQ_PASSWORD + "@" + config.RABBITMQ_HOST + ":" + config.RABBITMQ_PORT + "/")
	lib.FailOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	lib.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		queue,
		false,
		true,
		false,
		false,
		nil)
	lib.FailOnError(err, "Failed to declare a queue")

	err = ch.Qos(
		1,
		0,
		false,
	)
	lib.FailOnError(err, "Failed to set QoS")

	msgs, err := ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	lib.FailOnError(err, "Failed to register a consumer")

	for msg := range msgs {
		logrus.Println("Received a task")
		var msgBody rabbit.QueueMsg
		err = json.Unmarshal(msg.Body, &msgBody)
		if err != nil {
			logrus.Warn(err)
		}

		task := rabbitTask{
			msg:     msg,
			ch:      ch,
			msgBody: msgBody,
		}

		switch msgBody.TaskId {
		case rabbit.GAME_INSTALL:
			logrus.Info("INSTALL_GAME_SERVER")

			game := lib.GAMES[task.msgBody.Game]

			logrus.Info("Creating gs dir")

			gsPath := "/opt/aimpanel/gs/" + task.msgBody.GameServerID
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
					arg = strings.Replace(arg, "{uuid}", task.msgBody.GameServerID, -1)
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

			rabbitRpcSimpleResponse(task, rabbit.QueueMsg{
				TaskEnd: true,
				TaskOk:  true,
			})

			msg.Ack(false)
		case rabbit.WRAPPER_START:
			logrus.Info("START_WRAPPER")
			cmd := exec.Command("slave", "wrapper", task.msgBody.GameServerID)
			if err := cmd.Start(); err != nil {
				logrus.Error(err)
			}
			cmd.Process.Release()

			rabbitRpcSimpleResponse(task, rabbit.QueueMsg{
				TaskEnd: true,
				TaskOk:  true,
			})

			msg.Ack(false)
		}
	}
}

func rabbitRpcSimpleResponse(task rabbitTask, msg rabbit.QueueMsg) {
	body, err := json.Marshal(msg)
	err = task.ch.Publish(
		"",
		task.msg.ReplyTo,
		false,
		false,
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: task.msg.CorrelationId,
			Body:          body,
		})
	if err != nil {
		logrus.Error("Failed to respond: %v", err.Error())
		return
	}
}
