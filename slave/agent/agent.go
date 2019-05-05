package agent

import (
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"gitlab.com/systemz/aimpanel2/lib"
)

var (
	conn     *amqp.Connection
	channel  *amqp.Channel
	rpcQueue amqp.Queue
	err      error
)

func Start(token string) {
	log.Info("Starting Agent")

	// Defer can't be in init because this will be executed when the function return.
	conn, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
	lib.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	channel, err = conn.Channel()
	lib.FailOnError(err, "Failed to open channel")
	defer channel.Close()

	rpcQueue, err = channel.QueueDeclare("agent_rpc", false, false, false, false, nil)
	lib.FailOnError(err, "Failed to declare a rpc queue")

	err = channel.Qos(1, 0, false)
	lib.FailOnError(err, "Failed to set QoS")

	go Rpc(channel, rpcQueue)

	select {}
}
