package main_test

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"gitlab.com/systemz/aimpanel2/lib"
	"gitlab.com/systemz/aimpanel2/lib/rabbit"
	"testing"
)

func TestAgentInstallGameServer(t *testing.T) {
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

	start := rabbit.QueueMsg{
		TaskId:       rabbit.GAME_INSTALL,
		Game:         "minecraft",
		GameServerID: "test-test-test-test",
	}

	jsonMarshal, _ := json.Marshal(start)

	err = channel.Publish(
		"",
		"agent_1234",
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

func TestAgentStartWrapper(t *testing.T) {
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

	start := rabbit.QueueMsg{
		TaskId:       rabbit.WRAPPER_START,
		Game:         "minecraft",
		GameServerID: "test-test-test-test",
	}
	jsonMarshal, _ := json.Marshal(start)

	err = channel.Publish(
		"",
		"agent_1234",
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
