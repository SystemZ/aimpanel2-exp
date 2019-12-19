package wrapper

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"gitlab.com/systemz/aimpanel2/lib"
	"gitlab.com/systemz/aimpanel2/lib/rabbit"
	"gitlab.com/systemz/aimpanel2/slave/config"
	"net/http"
	"os"
)

var (
	channel *amqp.Channel
	queue   amqp.Queue
)

func Start(gameServerID string) {
	logrus.Info("Starting wrapper")

	token := os.Getenv("HOST_TOKEN")
	resp, err := http.Get(config.API_URL + "/v1/host/credentials/" + token + "/gs/" + gameServerID)
	if err != nil {
		lib.FailOnError(err, "Failed to get rabbit credentials")
	}

	var creds rabbit.Credentials
	err = json.NewDecoder(resp.Body).Decode(&creds)
	if err != nil {
		lib.FailOnError(err, "Failed to decode rabbit credentials json")
	}

	conn, err := amqp.Dial("amqp://" + creds.Username + ":" + creds.Password + "@" + creds.Host + ":" + creds.Port + creds.VHost)
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
	go p.Heartbeat()

	logrus.Info("Send WRAPPER_STARTED")
	p.SendToQueueData(rabbit.WRAPPER_STARTED)

	logrus.Info("Send WRAPPER_METRICS_FREQUENCY")
	p.SendToQueueData(rabbit.WRAPPER_METRICS_FREQUENCY)

	select {}
}
