package main

import (
	"aimpanel2/lib"
	"bufio"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"io"
	"os/exec"
)

var (
	channel  *amqp.Channel
	queue    amqp.Queue
	rpcQueue amqp.Queue
)

type Wrapper struct {
	Command string
	Args    string

	Output chan string
	Input  chan string
}

func (w *Wrapper) Run() {
	cmd := exec.Command("java", "-jar", "bungee/BungeeCord.jar")

	stdout, _ := cmd.StdoutPipe()
	stdin, _ := cmd.StdinPipe()

	if err := cmd.Start(); err != nil {
		log.Fatal("cmd.Start()", err)
	}

	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		w.Output <- scanner.Text()
	}

	for {
		log.Info("Got message to execute")
		in := <-w.Input
		io.WriteString(stdin, in+"\r\n")
	}

	if err := cmd.Wait(); err != nil {
		log.Fatal("cmd.Wait()", err)
	}
}

func (w *Wrapper) Log() {
	for {
		msg := <-w.Output

		err := channel.Publish("", queue.Name, false, false, amqp.Publishing{ContentType: "text/plain", Body: []byte(msg)})
		failOnError(err, "Publish error")

		log.WithFields(log.Fields{
			"method": "Collect",
			"msg":    msg,
		}).Info()
	}
}

func (w *Wrapper) Rpc() {
	msgs, err := channel.Consume(rpcQueue.Name, "", false, false, false, false, nil)
	failOnError(err, "Failed to register a consumer")

	for d := range msgs {
		log.Info("Got rpc message from rabbit")
		var wr lib.RpcMessage
		err := json.Unmarshal(d.Body, &wr)
		if err != nil {
			log.Warn(err)
		}

		switch wr.Type {
		case lib.START:
			log.Info("startServer message")
			go w.Run()

			err = channel.Publish("", d.ReplyTo, false, false, amqp.Publishing{
				ContentType:   "application/json",
				CorrelationId: d.CorrelationId,
				Body:          []byte(""),
			})
			failOnError(err, "Failed to publish a message")

			d.Ack(false)
		case lib.COMMAND:
			log.Info("sendMessage message")

			w.Input <- string(wr.Body)

			err = channel.Publish("", d.ReplyTo, false, false, amqp.Publishing{
				ContentType:   "application/json",
				CorrelationId: d.CorrelationId,
				Body:          []byte(""),
			})
			failOnError(err, "Failed to publish a message")

			d.Ack(false)
		case lib.STOP_SIGKILL:
			break
		case lib.STOP_SIGTERM:
			break
		}
	}
}

func main() {
	log.Info("Starting wrapper")

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

	output := make(chan string)
	input := make(chan string)

	w := &Wrapper{Output: output, Input: input}
	go w.Log()
	go w.Rpc()

	select {}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatal("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}
