package cron

import (
	"encoding/base64"
	"github.com/alexandrevicenzi/go-sse"
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/lib/task"
	"gitlab.com/systemz/aimpanel2/master/config"
	"gitlab.com/systemz/aimpanel2/master/events"
	"gitlab.com/systemz/aimpanel2/master/model"
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

func CheckHostsHeartbeat() {
	for {
		<-time.After(15 * time.Second)

		hosts, err := model.GetHosts()
		if err != nil {
			logrus.Error(err)
			return
		}

		for _, host := range hosts {
			lastTimestamp, err := model.Redis.Get("agent_heartbeat_token_" + host.Token).Int64()
			if err == nil {
				heartbeatTime := time.Unix(lastTimestamp, 0)

				if time.Since(heartbeatTime) > 10*time.Second {
					if host.State == 1 {
						host.State = 0
						err := model.Put(&host)
						if err != nil {
							logrus.Error(err)
						}
					}

				} else {
					if host.State == 0 {
						host.State = 1
						err := model.Put(&host)
						if err != nil {
							logrus.Error(err)
						}
					}

				}
			}
		}

	}
}

func CheckGSHeartbeat() {
	for {
		<-time.After(15 * time.Second)

		gameServers, err := model.GetGameServers()
		if err != nil {
			logrus.Error(err)
			return
		}

		for _, gs := range gameServers {
			lastTimestamp, err := model.Redis.Get("wrapper_heartbeat_id_" + gs.ID.Hex()).Int64()
			if err == nil {
				heartbeatTime := time.Unix(lastTimestamp, 0)

				if time.Since(heartbeatTime) > 10*time.Second {
					if gs.State == 1 {
						gs.State = 0
						err := model.Update(&gs)
						if err != nil {
							logrus.Error(err)
						}
					}

				} else {
					if gs.State == 0 {
						gs.State = 1
						err := model.Update(&gs)
						if err != nil {
							logrus.Error(err)
						}
					}

				}
			}
		}

	}
}
