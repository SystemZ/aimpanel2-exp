// +build integration

package main_test

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"gitlab.com/systemz/aimpanel2/lib"
	"testing"
)

func TestWrapperStart(t *testing.T) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		t.Fatal(err, "Failed to connect to RabbitMQ")
	}

	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		t.Fatal(err, "Failed to open a channel")
	}

	defer channel.Close()

	queue, err := channel.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil)
	if err != nil {
		t.Fatal(err, "Failed to declare a queue")
	}

	corrId := lib.RandomString(32)

	start := lib.RpcMessage{
		Type:           lib.START,
		Body:           "",
		Game:           "minecraft",
		GameServerUUID: "test-test-test-test",
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
	if err != nil {
		t.Fatal(err, "Failed to publish a message")
	}
}

func TestWrapperCommand(t *testing.T) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		t.Fatal(err, "Failed to connect to RabbitMQ")
	}

	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		t.Fatal(err, "Failed to open a channel")
	}

	defer channel.Close()

	queue, err := channel.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil)
	if err != nil {
		t.Fatal(err, "Failed to declare a queue")
	}

	corrId := lib.RandomString(32)

	start := lib.RpcMessage{
		Type:           lib.COMMAND,
		Body:           "alert test",
		Game:           "minecraft",
		GameServerUUID: "test-test-test-test",
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
	if err != nil {
		t.Fatal(err, "Failed to publish a message")
	}
}

func TestWrapperStopSigkill(t *testing.T) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		t.Fatal(err, "Failed to connect to RabbitMQ")
	}

	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		t.Fatal(err, "Failed to open a channel")
	}

	defer channel.Close()

	queue, err := channel.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil)
	if err != nil {
		t.Fatal(err, "Failed to declare a queue")
	}

	corrId := lib.RandomString(32)

	start := lib.RpcMessage{
		Type:           lib.STOP_SIGKILL,
		Body:           "",
		Game:           "minecraft",
		GameServerUUID: "test-test-test-test",
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
	if err != nil {
		t.Fatal(err, "Failed to publish a message")
	}
}

func TestWrapperStopSigterm(t *testing.T) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		t.Fatal(err, "Failed to connect to RabbitMQ")
	}

	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		t.Fatal(err, "Failed to open a channel")
	}

	defer channel.Close()

	queue, err := channel.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil)
	if err != nil {
		t.Fatal(err, "Failed to declare a queue")
	}

	corrId := lib.RandomString(32)

	start := lib.RpcMessage{
		Type:           lib.STOP_SIGTERM,
		Body:           "",
		Game:           "minecraft",
		GameServerUUID: "test-test-test-test",
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
	if err != nil {
		t.Fatal(err, "Failed to publish a message")
	}
}
