package agent

import (
	log "github.com/sirupsen/logrus"
)

func Start(token string) {
	log.Info("Starting Agent")

	rabbitListen("agent_" + token)

	select {}
}
