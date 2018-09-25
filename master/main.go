package main

import (
	"github.com/streadway/amqp"
	"log"
	"math/rand"
)

func main() {
	conn, err := amqp.Dial("amqp://admin:admin@46.105.209.74:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	channel, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer channel.Close()

	queue, err := channel.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil)
	failOnError(err, "Failed to declare a queue")

	//msgs, err := channel.Consume(
	//	queue.Name,
	//	"",
	//	true,
	//	false,
	//	false,
	//	false,
	//	nil)
	//failOnError(err, "Failed to register a consumer")

	corrId := randomString(32)

	err = channel.Publish(
		"",
		"wrapper_rpc",
		false,
		false,
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrId,
			ReplyTo:       queue.Name,
			Type:          "sendMessage",
			Body:          []byte("Hello!"),
		})

	failOnError(err, "Failed to publish a message")
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}
