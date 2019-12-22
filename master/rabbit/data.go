package rabbit

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/lib"
	"gitlab.com/systemz/aimpanel2/lib/rabbit"
	"gitlab.com/systemz/aimpanel2/master/model"
	"time"
)

func ListenAgentData() {
	queue, err := channel.QueueDeclare(
		"agent_data",
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	lib.FailOnError(err, "Failed to declare a queue")

	msgs, err := channel.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil)
	lib.FailOnError(err, "Failed to register a consumer")

	go func() {
		for msg := range msgs {
			var msgBody rabbit.QueueMsg
			err = json.Unmarshal(msg.Body, &msgBody)
			if err != nil {
				logrus.Warn(err)
			}

			switch msgBody.TaskId {
			case rabbit.AGENT_METRICS_FREQUENCY:
				agentToken := msgBody.AgentToken

				var host model.Host
				if model.DB.Where("token = ?", agentToken).First(&host).RecordNotFound() {
					break
				}

				msg := rabbit.QueueMsg{
					TaskId:          rabbit.AGENT_METRICS_FREQUENCY,
					MetricFrequency: host.MetricFrequency,
				}

				err := SendRpcMessage("agent_"+agentToken, msg)
				if err != nil {
					logrus.Error(err.Error())
				}
			case rabbit.AGENT_METRICS:
				agentToken := msgBody.AgentToken

				var host model.Host
				if model.DB.Where("token = ?", agentToken).First(&host).RecordNotFound() {
					break
				}

				metric := &model.MetricHost{
					HostId:    host.ID,
					CpuUsage:  msgBody.CpuUsage,
					RamFree:   msgBody.RamFree,
					RamTotal:  msgBody.RamTotal,
					DiskFree:  msgBody.DiskFree,
					DiskUsed:  msgBody.DiskUsed,
					DiskTotal: msgBody.DiskTotal,
					User:      msgBody.User,
					System:    msgBody.System,
					Idle:      msgBody.Idle,
					Nice:      msgBody.Nice,
					Iowait:    msgBody.Iowait,
					Irq:       msgBody.Irq,
					Softirq:   msgBody.Softirq,
					Steal:     msgBody.Steal,
					Guest:     msgBody.Guest,
					GuestNice: msgBody.GuestNice,
				}
				model.DB.Save(metric)
			case rabbit.AGENT_OS:
				agentToken := msgBody.AgentToken

				var host model.Host
				if model.DB.Where("token = ?", agentToken).First(&host).RecordNotFound() {
					break
				}

				host.OS = msgBody.OS
				host.Platform = msgBody.Platform
				host.PlatformFamily = msgBody.PlatformFamily
				host.PlatformVersion = msgBody.PlatformVersion
				host.KernelVersion = msgBody.KernelVersion
				host.KernelArch = msgBody.KernelArch

				model.DB.Save(&host)

			case rabbit.AGENT_HEARTBEAT:
				agentToken := msgBody.AgentToken
				model.Redis.Set("agent_heartbeat_token_"+agentToken, msgBody.Timestamp, 24*time.Hour)
			}
		}
	}()
}
