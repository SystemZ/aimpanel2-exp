package main

import (
	"aimpanel2/lib"
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"io"
	"log"
	"os/exec"
)

var (
	channel  *amqp.Channel
	queue    amqp.Queue
	rpcQueue amqp.Queue

	stdin  io.WriteCloser
	stdout io.ReadCloser
)

func main() {
	conn, err := amqp.Dial("amqp://admin:admin@46.105.209.74:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	channel, err = conn.Channel()
	failOnError(err, "Failed to open channel")
	defer channel.Close()

	queue, err = channel.QueueDeclare("wrapper", false, false, false, false, nil)
	failOnError(err, "Failed to declare a queue")

	rpcQueue, err = channel.QueueDeclare("wrapper_rpc", false, false, false, false, nil)
	failOnError(err, "Failed to declare a rpc queue")

	err = channel.Qos(1, 0, false)
	failOnError(err, "Failed to set QoS")

	//RPC
	msgs, err := channel.Consume(rpcQueue.Name, "", false, false, false, false, nil)
	failOnError(err, "Failed to register a consumer")

	go func() {
		for d := range msgs {
			var wr lib.WrapperRPC
			err := json.Unmarshal(d.Body, &wr)
			if err != nil {
				log.Println(err)
			}

			switch wr.Type {
			case lib.START:
				log.Println("startServer")
				startServer()
			case lib.COMMAND:
				log.Println("sendMessage")
				io.WriteString(stdin, string(d.Body)+"\r\n")

				err = channel.Publish("", d.ReplyTo, false, false, amqp.Publishing{
					ContentType:   "application/json",
					CorrelationId: d.CorrelationId,
					Body:          []byte(""),
				})
				failOnError(err, "Failed to publish a message")

				d.Ack(false)
			}
		}
	}()

	select {}
}

func startServer() {
	cmd := exec.Command("bash", "fake-server.sh")

	stdout, _ = cmd.StdoutPipe()
	stdin, _ = cmd.StdinPipe()

	if err := cmd.Start(); err != nil {
		log.Fatalln(err)
	}

	readStdout()

	if err := cmd.Wait(); err != nil {
		log.Println(err)
	}
}

func readStdout() {
	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		err := channel.Publish("", queue.Name, false, false, amqp.Publishing{ContentType: "text/plain", Body: []byte(scanner.Text())})
		failOnError(err, "Publish error")
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}
