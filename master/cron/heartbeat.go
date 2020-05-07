package cron

import (
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/master/model"
	"time"
)

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
