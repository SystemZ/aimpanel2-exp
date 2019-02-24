package rabbit

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"gitlab.com/systemz/aimpanel2/lib"
	"gitlab.com/systemz/aimpanel2/lib/rabbit"
	"gitlab.com/systemz/aimpanel2/slave/config"
)

var (
	channel *amqp.Channel
)

func RabbitListen() {
	conn, err := amqp.Dial("amqp://" + config.RABBITMQ_USERNAME + ":" + config.RABBITMQ_PASSWORD + "@" + config.RABBITMQ_HOST + ":" + config.RABBITMQ_PORT + "/")
	lib.FailOnError(err, "Failed to connect to RabbitMQ")

	channel, err = conn.Channel()
	lib.FailOnError(err, "Failed to open a channel")
	defer channel.Close()

	err = channel.Qos(
		1,
		0,
		false,
	)
	lib.FailOnError(err, "Failed to set QoS")
}

func SendRpcMessage(queue string, msg rabbit.QueueMsg) {
	body, err := json.Marshal(msg)

	corrId := lib.RandomString(32)

	err = channel.Publish(
		"",
		queue,
		false,
		false,
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: corrId,
			Body:          body,
		})
	if err != nil {
		logrus.Error("Failed to respond: %v", err.Error())
		return
	}
}
