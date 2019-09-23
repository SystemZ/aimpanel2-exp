package cron

import (
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/master/model"
	"time"
)

func CheckHeartbeat() {
	for {
		<-time.After(15 * time.Second)

		var hosts []model.Host
		model.DB.Find(&hosts)

		for _, host := range hosts {
			lastTimestamp, err := model.Redis.Get("agent_heartbeat_token_" + host.Token).Int64()
			if err == nil {
				heartbeatTime := time.Unix(lastTimestamp, 0)

				logrus.Info(host.ID)
				logrus.Info(time.Since(heartbeatTime))
			}
		}

	}
}
