package main

import (
	"bufio"
	"fmt"
	"github.com/streadway/amqp"
	"io"
	"log"
	"os/exec"
	"time"
)

var (
	ch *amqp.Channel
	q  amqp.Queue

	stdin  io.WriteCloser
	stdout io.ReadCloser
)

func main() {
	conn, err := amqp.Dial("amqp://admin:admin@46.105.209.74:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err = conn.Channel()
	failOnError(err, "Failed to open channel")
	defer ch.Close()

	q, err = ch.QueueDeclare("wrapper", false, false, false, false, nil)
	failOnError(err, "Failed to declare a queue")

	startServer()
}

func startServer() {
	cmd := exec.Command("bash", "fake-server.sh")

	stdout, _ = cmd.StdoutPipe()
	stdin, _ = cmd.StdinPipe()

	if err := cmd.Start(); err != nil {
		log.Fatalln(err)
	}

	go func() {
		for {
			io.WriteString(stdin, "test\r\n")
			time.Sleep(4 * time.Second)
		}
	}()

	readStdout()

	if err := cmd.Wait(); err != nil {
		log.Println(err)
	}
}

func readStdout() {
	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		log.Println(scanner.Text())
		err := ch.Publish("", q.Name, false, false, amqp.Publishing{ContentType: "text/plain", Body: []byte(scanner.Text())})
		failOnError(err, "Publish error")
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}
