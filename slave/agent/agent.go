package agent

import (
	"bytes"
	"encoding/json"
	"github.com/r3labs/sse"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"gitlab.com/systemz/aimpanel2/lib"
	"gitlab.com/systemz/aimpanel2/lib/rabbit"
	"gitlab.com/systemz/aimpanel2/lib/response"
	"gitlab.com/systemz/aimpanel2/slave/config"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"
)

var (
	hostToken        string
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
	hostToken = t

	resp, err := http.Get(config.API_URL + "/v1/host/credentials/" + hostToken)
	if err != nil {
		lib.FailOnError(err, "Failed to get rabbit credentials")
	}

	var creds rabbit.Credentials
	err = json.NewDecoder(resp.Body).Decode(&creds)
	if err != nil {
		lib.FailOnError(err, "Failed to decode rabbit credentials json")
	}

	resp, err = http.Get(config.API_URL + "/v1/host/auth/" + hostToken)
	if err != nil {
		lib.FailOnError(err, "Failed to get host token")
	}

	var token response.Token
	err = json.NewDecoder(resp.Body).Decode(&token)
	if err != nil {
		lib.FailOnError(err, "Failed to decode rabbit credentials json")
	}

	conn, err := amqp.Dial("amqp://" + creds.Username + ":" + creds.Password + "@" + creds.Host + ":" + creds.Port + creds.VHost)
	lib.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	channel, err = conn.Channel()
	lib.FailOnError(err, "Failed to open a channel")
	defer channel.Close()

	queue, err = channel.QueueDeclare(
		"agent_"+hostToken,
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

	go agent(token.Token)

	//sendToQueueData(rabbit.AGENT_METRICS_FREQUENCY)

	//sendOSInfo()

	//go heartbeat()

	select {}
}

func agent(token string) {
	client := sse.NewClient(config.API_URL + "/v1/events/" + hostToken)
	client.Headers = map[string]string{
		"Authorization": "Bearer " + token,
	}
	err := client.SubscribeRaw(func(msg *sse.Event) {
		logrus.Info(string(msg.ID))
		logrus.Info(string(msg.Data))
		logrus.Info(string(msg.Event))

		gsId := string(msg.Data)
		taskId, _ := strconv.Atoi(string(msg.Event))
		switch taskId {
		case rabbit.WRAPPER_START:
			logrus.Info("START_WRAPPER")
			cmd := exec.Command("slave", "wrapper", gsId)

			cmd.Env = os.Environ()
			cmd.Env = append(cmd.Env, "HOST_TOKEN="+hostToken)
			cmd.Env = append(cmd.Env, "API_TOKEN="+token)

			//TODO: FOR TESTING ONLY
			var stdBuffer bytes.Buffer
			mw := io.MultiWriter(os.Stdout, &stdBuffer)
			cmd.Stdout = mw
			cmd.Stderr = mw

			if err := cmd.Start(); err != nil {
				logrus.Error(err)
			}

			cmd.Process.Release()
		}
	})
	if err != nil {
		lib.FailOnError(err, "Failed to subscribe a channel")
	}

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

			gsPath := filepath.Clean(config.GS_DIR) + "/" + task.msgBody.GameServerID.String()
			if _, err := os.Stat(gsPath); os.IsNotExist(err) {
				//TODO: Set correct perms
				_ = os.Mkdir(gsPath, 0777)
			}

			err = task.msgBody.Game.Install(filepath.Clean(config.STORAGE_DIR), gsPath)
			if err != nil {
				logrus.Error(err)
			}

			logrus.Info("Installation finished")

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
		case rabbit.AGENT_REMOVE_GS:
			logrus.Info("AGENT_REMOVE_GS")

			gsId := task.msgBody.GameServerID.String()
			gsPath := filepath.Clean(config.GS_DIR) + "/" + gsId
			gsTrashPath := filepath.Clean(config.TRASH_DIR) + "/" + gsId

			err := os.Rename(gsPath, gsTrashPath)
			if err != nil {
				logrus.Error(err)
			}

			msg.Ack(false)
		}
	}
}

func metrics() {
	for {
		<-time.After(time.Duration(metricsFrequency) * time.Second)

		virtualMemory, _ := mem.VirtualMemory()
		ramFree := virtualMemory.Free / 1024 / 1024
		ramTotal := virtualMemory.Total / 1024 / 1024
		cpuPercent, _ := cpu.Percent(time.Duration(1)*time.Second, false)
		cpuTimes, _ := cpu.Times(false)

		diskUsage, _ := disk.Usage("/")
		diskFree := diskUsage.Free / 1024 / 1024
		diskTotal := diskUsage.Total / 1024 / 1024
		diskUsed := diskUsage.Used / 1024 / 1024

		msg := rabbit.QueueMsg{
			TaskId:     rabbit.AGENT_METRICS,
			AgentToken: hostToken,
			CpuUsage:   int(cpuPercent[0]),
			RamFree:    int(ramFree),
			RamTotal:   int(ramTotal),
			DiskFree:   int(diskFree),
			DiskTotal:  int(diskTotal),
			DiskUsed:   int(diskUsed),
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

func heartbeat() {
	for {
		<-time.After(5 * time.Second)

		logrus.Info("Sending heartbeat")

		msg := rabbit.QueueMsg{
			TaskId:     rabbit.AGENT_HEARTBEAT,
			AgentToken: hostToken,
			Timestamp:  time.Now().Unix(),
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
		AgentToken: hostToken,
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

func sendOSInfo() {
	h, _ := host.Info()

	msg := rabbit.QueueMsg{
		TaskId:     rabbit.AGENT_OS,
		AgentToken: hostToken,

		OS:              h.OS,
		Platform:        h.Platform,
		PlatformFamily:  h.PlatformFamily,
		PlatformVersion: h.PlatformVersion,
		KernelVersion:   h.KernelVersion,
		KernelArch:      h.KernelArch,
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
