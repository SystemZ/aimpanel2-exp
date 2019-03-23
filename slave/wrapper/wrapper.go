package wrapper

import (
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"gitlab.com/systemz/aimpanel2/lib"
	"gitlab.com/systemz/aimpanel2/slave/config"
)

var (
	conn     *amqp.Connection
	channel  *amqp.Channel
	queue    amqp.Queue
	rpcQueue amqp.Queue
	err      error
)

func Start(gameServerID string) {
	logrus.Info("Starting wrapper")
	//TODO: Make request to master to get creds to rabbit

	// Defer can't be in init because this will be executed when the function return.
	conn, err := amqp.Dial("amqp://" + config.RABBITMQ_USERNAME + ":" + config.RABBITMQ_PASSWORD + "@" + config.RABBITMQ_HOST + ":" + config.RABBITMQ_PORT + "/")
	lib.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	channel, err = conn.Channel()
	lib.FailOnError(err, "Failed to open channel")
	defer channel.Close()

	queue, err = channel.QueueDeclare("wrapper_"+gameServerID, true, false, false, false, nil)
	lib.FailOnError(err, "Failed to declare a low queue")

	err = channel.Qos(1, 0, false)
	lib.FailOnError(err, "Failed to set QoS")

	output := make(chan string)
	input := make(chan string)

	p := &Process{
		Output: output,
		Input:  input,

		//amqp
		Channel: channel,
		Queue:   queue,

		GameServerID: gameServerID,
	}

	go p.LogStdout()
	go p.LogStderr()
	go p.Rpc()

	//go p.Run()

	select {}
}
