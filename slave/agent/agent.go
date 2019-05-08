package agent

import (
	"encoding/json"
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"gitlab.com/systemz/aimpanel2/lib"
	"gitlab.com/systemz/aimpanel2/lib/rabbit"
	"gitlab.com/systemz/aimpanel2/slave/config"
	"os"
	"os/exec"
	"strings"
	"time"
)

var (
	token            string
	channel          *amqp.Channel
	queue            amqp.Queue
	queueData        amqp.Queue
	metricsFrequency int
)

type rabbitTask struct {
	msg     amqp.Delivery
	msgBody rabbit.QueueMsg
	ch      *amqp.Channel
}

func Start(t string) {
	logrus.Info("Starting Agent")
	token = t

	metrics()

	conn, err := amqp.Dial("amqp://" + config.RABBITMQ_USERNAME + ":" + config.RABBITMQ_PASSWORD + "@" + config.RABBITMQ_HOST + ":" + config.RABBITMQ_PORT + "/")
	lib.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	channel, err = conn.Channel()
	lib.FailOnError(err, "Failed to open a channel")
	defer channel.Close()

	queue, err = channel.QueueDeclare(
		"agent_"+token,
		true,
		false,
		false,
		false,
		nil)
	lib.FailOnError(err, "Failed to declare a queue")

	queueData, err = channel.QueueDeclare(
		"agent_data",
		true,
		false,
		false,
		false,
		nil)
	lib.FailOnError(err, "Failed to declare a queue")

	err = channel.Qos(
		1,
		0,
		false,
	)
	lib.FailOnError(err, "Failed to set QoS")

	go agent()

	sendToQueueData(rabbit.AGENT_METRICS_FREQUENCY)

	select {}
}

func agent() {
	msgs, err := channel.Consume(
		queue.Name,
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
			ch:      channel,
			msgBody: msgBody,
		}

		switch msgBody.TaskId {
		case rabbit.GAME_INSTALL:
			logrus.Info("INSTALL_GAME_SERVER")

			logrus.Info("Creating gs dir")

			gameFile := task.msgBody.GameFile

			gsPath := config.GS_DIR + task.msgBody.GameServerID.String()
			if _, err := os.Stat(gsPath); os.IsNotExist(err) {
				os.Mkdir(gsPath, 0777)
			}

			logrus.Info("Downloading install package")

			fileNameWithVersion := fmt.Sprintf("%d_%s", gameFile.GameVersion, gameFile.Filename)

			if _, err = os.Stat(config.STORAGE_DIR + fileNameWithVersion); os.IsNotExist(err) {
				cmd := exec.Command("wget", "-O", fileNameWithVersion, gameFile.DownloadUrl)
				cmd.Dir = config.STORAGE_DIR

				if err := cmd.Run(); err != nil {
					logrus.Error(err)
				}

				cmd.Wait()
			}

			logrus.Info("Executing install commands")

			for _, c := range *task.msgBody.GameCommands {
				c.Command = strings.Replace(c.Command, "{storageDir}", config.STORAGE_DIR, -1)
				c.Command = strings.Replace(c.Command, "{gsDir}", config.GS_DIR, -1)
				c.Command = strings.Replace(c.Command, "{uuid}", task.msgBody.GameServerID.String(), -1)
				c.Command = strings.Replace(c.Command, "{fileName}", fileNameWithVersion, -1)

				command := strings.Split(c.Command, " ")

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
			cmd := exec.Command("slave", "wrapper", task.msgBody.GameServerID.String())
			if err := cmd.Start(); err != nil {
				logrus.Error(err)
			}
			cmd.Process.Release()

			rabbitRpcSimpleResponse(task, rabbit.QueueMsg{
				TaskEnd: true,
				TaskOk:  true,
			})

			msg.Ack(false)

		case rabbit.AGENT_METRICS_FREQUENCY:
			logrus.Info("AGENT_METRICS_FREQUENCY")

			metricsFrequency = msgBody.MetricFrequency

			go metrics()

			rabbitRpcSimpleResponse(task, rabbit.QueueMsg{
				TaskEnd: true,
				TaskOk:  true,
			})

			msg.Ack(false)
		}
	}
}

func metrics() {
	for {
		<-time.After(time.Duration(metricsFrequency) * time.Second)

		virtualMemory, _ := mem.VirtualMemory()
		ramFree := virtualMemory.Free / 1024 / 1024
		cpuPercent, _ := cpu.Percent(time.Duration(1)*time.Second, false)
		cpuTimes, _ := cpu.Times(false)

		diskUsage, _ := disk.Usage("/")

		msg := rabbit.QueueMsg{
			TaskId:     rabbit.AGENT_METRICS,
			AgentToken: token,
			CpuUsage:   int(cpuPercent[0]),
			RamFree:    int(ramFree),
			DiskFree:   int(diskUsage.Free),
			User:       int(cpuTimes[0].User),
			System:     int(cpuTimes[0].System),
			Idle:       int(cpuTimes[0].Idle),
			Nice:       int(cpuTimes[0].Nice),
			Iowait:     int(cpuTimes[0].Iowait),
			Irq:        int(cpuTimes[0].Irq),
			Softirq:    int(cpuTimes[0].Softirq),
			Steal:      int(cpuTimes[0].Steal),
			Guest:      int(cpuTimes[0].Guest),
			GuestNice:  int(cpuTimes[0].GuestNice),
			Stolen:     int(cpuTimes[0].Stolen),
		}

		msgJson, _ := json.Marshal(msg)

		err := channel.Publish(
			"",
			queueData.Name,
			false,
			false,
			amqp.Publishing{
				ContentType:   "application/json",
				CorrelationId: lib.RandomString(32),
				Body:          msgJson,
			})

		lib.FailOnError(err, "Publish error")
	}
}

func sendToQueueData(taskId int) {
	msg := rabbit.QueueMsg{
		TaskId:     taskId,
		AgentToken: token,
	}

	msgJson, _ := json.Marshal(msg)

	err := channel.Publish(
		"",
		queueData.Name,
		false,
		false,
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: lib.RandomString(32),
			Body:          msgJson,
		})

	lib.FailOnError(err, "Publish error")
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
		logrus.Errorf("Failed to respond: %v", err.Error())
		return
	}
}
