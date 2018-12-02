package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"gitlab.com/systemz/aimpanel2/lib"
	"gitlab.com/systemz/aimpanel2/wrapper/process"
)

var (
	conn        *amqp.Connection
	channel     *amqp.Channel
	queueLow    amqp.Queue
	queueNormal amqp.Queue
	queueHigh   amqp.Queue
	rpcQueue    amqp.Queue
	err         error
)

func init() {
	log.Info("Init wrapper")

	// Defer can't be in init because this will be executed when the function return.

	conn, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
	lib.FailOnError(err, "Failed to connect to RabbitMQ")

	channel, err = conn.Channel()
	lib.FailOnError(err, "Failed to open channel")

	queueLow, err = channel.QueueDeclare("wrapper_low", true, false, false, false, nil)
	lib.FailOnError(err, "Failed to declare a low queue")

	queueNormal, err = channel.QueueDeclare("wrapper_normal", true, false, false, false, nil)
	lib.FailOnError(err, "Failed to declare a normal queue")

	queueHigh, err = channel.QueueDeclare("wrapper_high", true, false, false, false, nil)
	lib.FailOnError(err, "Failed to declare a high queue")

	rpcQueue, err = channel.QueueDeclare("wrapper_rpc", false, false, false, false, nil)
	lib.FailOnError(err, "Failed to declare a rpc queue")

	err = channel.Qos(1, 0, false)
	lib.FailOnError(err, "Failed to set QoS")
}

func main() {
	log.Info("Starting wrapper")

	defer conn.Close()
	defer channel.Close()

	output := make(chan string)
	input := make(chan string)

	p := &process.Process{
		Output: output,
		Input:  input,

		//amqp
		Channel:     channel,
		QueueLow:    queueLow,
		QueueNormal: queueNormal,
		QueueHigh:   queueHigh,
		RpcQueue:    rpcQueue,
	}
	go p.Log()
	go p.Rpc()

	select {}
}
