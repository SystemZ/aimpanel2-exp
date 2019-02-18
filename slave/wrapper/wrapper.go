package wrapper

import (
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"gitlab.com/systemz/aimpanel2/lib"
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

func Start(startToken string) {
	logrus.Info("Starting wrapper")
	//TODO: Make request to master to get creds to rabbit

	// Defer can't be in init because this will be executed when the function return.
	conn, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
	lib.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	channel, err = conn.Channel()
	lib.FailOnError(err, "Failed to open channel")
	defer channel.Close()

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

	output := make(chan string)
	input := make(chan string)

	p := &Process{
		Output: output,
		Input:  input,

		//amqp
		Channel:     channel,
		QueueLow:    queueLow,
		QueueNormal: queueNormal,
		QueueHigh:   queueHigh,
		RpcQueue:    rpcQueue,
	}

	go p.LogStdout()
	go p.LogStderr()
	go p.Rpc()

	select {}
}
