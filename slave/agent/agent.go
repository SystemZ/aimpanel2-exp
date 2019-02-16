package agent

import (
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

var (
	conn     *amqp.Connection
	channel  *amqp.Channel
	rpcQueue amqp.Queue
	err      error
)

func Start(token string) {
	log.Info("Starting Agent")

	rabbitListen("agent_" + token)

	select {}
}
