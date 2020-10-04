package cron

import (
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/master/config"
	"gitlab.com/systemz/aimpanel2/master/model"
	"time"
)

func RemoveOldMetrics() {
	go func() {
		for {
			time.Sleep(time.Hour * 1)

			hosts, err := model.GetHosts()
			if err != nil {
				logrus.Warn(err)
				continue
			}

			for _, host := range hosts {
				deleted, err := model.RemoveMetricsOlderThan(host.ID, model.HostMetric, host.MetricMaxS)
				if err != nil {
					logrus.Warn(err)
					continue
				}

				if config.DEV_MODE {
					logrus.Infof("Deleted %d host (%s) metrics", deleted, host.ID.Hex())
				}
			}

			gameServers, err := model.GetGameServers()
			if err != nil {
				logrus.Warn(err)
				continue
			}

			for _, gs := range gameServers {
				deleted, err := model.RemoveMetricsOlderThan(gs.ID, model.GameServerMetric, gs.MetricMaxS)
				if err != nil {
					logrus.Warn(err)
					continue
				}

				if config.DEV_MODE {
					logrus.Infof("Deleted %d gs (%s) metrics", deleted, gs.ID.Hex())
				}
			}
		}
	}()
}
