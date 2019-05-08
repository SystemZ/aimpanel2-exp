package wrapper

import (
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"gitlab.com/systemz/aimpanel2/lib"
	"gitlab.com/systemz/aimpanel2/lib/rabbit"
	"gitlab.com/systemz/aimpanel2/slave/config"
)

var (
	channel *amqp.Channel
	queue   amqp.Queue
)

func Start(gameServerID string) {
	logrus.Info("Starting wrapper")
	//TODO: Make request to master to get creds to rabbit

	// Defer can't be in init because this will be executed when the function return.
	conn, err := amqp.Dial("amqp://" + config.RABBITMQ_USERNAME + ":" + config.RABBITMQ_PASSWORD + "@" + config.RABBITMQ_HOST + ":" + config.RABBITMQ_PORT + config.RABBITMQ_VHOST)
	lib.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	channel, err = conn.Channel()
	lib.FailOnError(err, "Failed to open channel")
	defer channel.Close()

	queue, err = channel.QueueDeclare("wrapper_"+gameServerID, true,
		false, false, false, nil)
	lib.FailOnError(err, "Failed to declare a wrapper queue")

	queueData, err := channel.QueueDeclare("wrapper_data", true,
		false, false, false, nil)
	lib.FailOnError(err, "Failed to declare a data queue")

	err = channel.Qos(1, 0, false)
	lib.FailOnError(err, "Failed to set QoS")

	output := make(chan string)
	input := make(chan string)

	p := &Process{
		Output: output,
		Input:  input,

		//amqp
		Channel:   channel,
		Queue:     queue,
		QueueData: queueData,

		GameServerID: gameServerID,
	}

	go p.Rpc()

	logrus.Info("Send WRAPPER_STARTED")
	p.SendToQueueData(rabbit.WRAPPER_STARTED)

	logrus.Info("Send WRAPPER_METRICS_FREQUENCY")
	p.SendToQueueData(rabbit.WRAPPER_METRICS_FREQUENCY)

	select {}
}
