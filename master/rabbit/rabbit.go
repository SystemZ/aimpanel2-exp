package rabbit

import (
	"encoding/json"
	"github.com/michaelklishin/rabbit-hole"
	"github.com/streadway/amqp"
	"gitlab.com/systemz/aimpanel2/lib"
	rabbitLib "gitlab.com/systemz/aimpanel2/lib/rabbit"
	"gitlab.com/systemz/aimpanel2/master/config"
	"net/http"
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

func SetupRabbitAPI() {
	scheme := "http://"
	if config.RABBITMQ_TLS {
		scheme = "https://"
	}

	client, err := rabbithole.NewClient(scheme+config.RABBITMQ_HOST+":"+config.RABBITMQ_PORT_API, config.RABBITMQ_USERNAME, config.RABBITMQ_PASSWORD)
	if err != nil {
		lib.FailOnError(err, "Failed to connect to Rabbit API")
	}

	Client = client
}

func SendRpcMessage(queue string, msg rabbitLib.QueueMsg) error {
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

func PutUser(credentials rabbitLib.Credentials) error {
	resp, err := Client.PutUser(credentials.Username, rabbithole.UserSettings{Password: credentials.Password})
	if err != nil {
		return &lib.Error{ErrorCode: 1000}
	}

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		return &lib.Error{ErrorCode: 1001}
	}

	resp, err = Client.UpdatePermissionsIn(credentials.VHost, credentials.Username, rabbithole.Permissions{
		Configure: ".*",
		Write:     ".*",
		Read:      ".*",
	})
	if err != nil {
		return &lib.Error{ErrorCode: 1002}
	}

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		return &lib.Error{ErrorCode: 1003}
	}

	return nil
}
