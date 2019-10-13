package cron

import (
	"gitlab.com/systemz/aimpanel2/master/model"
	"time"
)

func CheckHostsHeartbeat() {
	for {
		<-time.After(15 * time.Second)

		var hosts []model.Host
		model.DB.Find(&hosts)

		for _, host := range hosts {
			lastTimestamp, err := model.Redis.Get("agent_heartbeat_token_" + host.Token).Int64()
			if err == nil {
				heartbeatTime := time.Unix(lastTimestamp, 0)

				if time.Since(heartbeatTime) > 10*time.Second {
					if host.State == 1 {
						host.State = 0
						model.DB.Save(&host)
					}

				} else {
					if host.State == 0 {
						host.State = 1
						model.DB.Save(&host)
					}

				}
			}
		}

	}
}

func CheckGSHeartbeat() {
	for {
		<-time.After(15 * time.Second)

		var gss []model.GameServer
		model.DB.Find(&gss)

		for _, gs := range gss {
			lastTimestamp, err := model.Redis.Get("wrapper_heartbeat_id_" + gs.ID.String()).Int64()
			if err == nil {
				heartbeatTime := time.Unix(lastTimestamp, 0)

				if time.Since(heartbeatTime) > 10*time.Second {
					if gs.State == 1 {
						gs.State = 0
						model.DB.Save(&gs)
					}

				} else {
					if gs.State == 0 {
						gs.State = 1
						model.DB.Save(&gs)
					}

				}
			}
		}

	}
}
