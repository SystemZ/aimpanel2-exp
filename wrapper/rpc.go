package main

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"gitlab.com/systemz/aimpanel2/lib"
	"log"
)

func main() {
	log.Println("start")
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	lib.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	channel, err := conn.Channel()
	lib.FailOnError(err, "Failed to open a channel")
	defer channel.Close()

	queue, err := channel.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil)
	lib.FailOnError(err, "Failed to declare a queue")

	//msgs, err := channel.Consume(
	//	queue.Name,
	//	"",
	//	true,
	//	false,
	//	false,
	//	false,
	//	nil)
	//lib.FailOnError(err, "Failed to register a consumer")

	corrId := lib.RandomString(32)

	start := lib.RpcMessage{
		Type: lib.STOP_SIGTERM,
		Body: "",
	}
	jsonMarshal, _ := json.Marshal(start)

	err = channel.Publish(
		"",
		"wrapper_rpc",
		false,
		false,
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: corrId,
			ReplyTo:       queue.Name,
			Body:          jsonMarshal,
		})

	lib.FailOnError(err, "Failed to publish a message")
}
