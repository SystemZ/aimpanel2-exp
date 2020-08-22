package cron

import (
	"encoding/base64"
	"github.com/alexandrevicenzi/go-sse"
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/lib/task"
	"gitlab.com/systemz/aimpanel2/master/config"
	"gitlab.com/systemz/aimpanel2/master/events"
	"strings"
	"time"
)

func SseHeartbeat() {
	// ready to send heartbeat messages
	raw1SlaveMsg := task.Message{TaskId: task.PING}
	raw2SlaveMsg, _ := raw1SlaveMsg.Serialize()
	slaveMsg := sse.SimpleMessage(raw2SlaveMsg)
	browserMsg := sse.SimpleMessage(base64.StdEncoding.EncodeToString([]byte("")))

	// send empty message as keep alive measure for SSE sessions
	go func() {
		for {
			time.Sleep(time.Second * 45)
			channels := events.SSE.Channels()
			// TODO we can check which hosts aren't connected and trigger alarm
			// send ping to all clients connected via SSE
			for _, ch := range channels {
				if strings.HasPrefix(ch, "/v1/events/") {
					// ping for agent
					events.SSE.SendMessage(
						ch,
						slaveMsg,
					)
				} else {
					// ping for browser
					events.SSE.SendMessage(
						ch,
						browserMsg,
					)
				}
			}
		}
	}()

	// show active connections in DEV env
	if config.DEV_MODE {
		go func() {
			for {
				channels := events.SSE.Channels()
				logrus.Infof("SSE channels: %v", channels)
				time.Sleep(time.Second * 30)
			}
		}()
	}
}
