package rabbit

import (
	"encoding/json"
	"github.com/michaelklishin/rabbit-hole"
	"github.com/streadway/amqp"
	"gitlab.com/systemz/aimpanel2/lib"
	"gitlab.com/systemz/aimpanel2/lib/rabbit"
	"gitlab.com/systemz/aimpanel2/master/config"
)

var (
	channel *amqp.Channel
	Client  *rabbithole.Client
)

func Listen() {
	conn, err := amqp.Dial("amqp://" + config.RABBITMQ_USERNAME + ":" + config.RABBITMQ_PASSWORD + "@" + config.RABBITMQ_HOST + ":" + config.RABBITMQ_PORT + config.RABBITMQ_VHOST)
	lib.FailOnError(err, "Failed to connect to RabbitMQ")

	channel, err = conn.Channel()
	lib.FailOnError(err, "Failed to open a channel")
	//defer channel.Close()

	err = channel.Qos(
		1,
		0,
		false,
	)
	lib.FailOnError(err, "Failed to set QoS")
}

func SetupRabbitHole() {
	//client, err := rabbithole.NewClient("http://" + config.RABBITMQ_HOST + ":" + config.RABBITMQ_PORT, config.RABBITMQ_USERNAME, config.RABBITMQ_PASSWORD)
	client, err := rabbithole.NewClient("http://localhost:15672", "guest", "guest")
	if err != nil {
		lib.FailOnError(err, "Failed to connect to Rabbit API")
	}

	Client = client
}

func SendRpcMessage(queue string, msg rabbit.QueueMsg) error {
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
		return err
	}

	return nil
}
