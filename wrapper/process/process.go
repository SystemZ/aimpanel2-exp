package process

import (
	"bufio"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"gitlab.com/systemz/aimpanel2/lib"
	"io"
	"os"
	"os/exec"
	"sync"
	"syscall"
)

type Process struct {
	Cmd *exec.Cmd

	Running bool

	Output chan string
	Input  chan string

	//amqp
	Channel     *amqp.Channel
	QueueLow    amqp.Queue
	QueueNormal amqp.Queue
	QueueHigh   amqp.Queue
	RpcQueue    amqp.Queue
}

func (p *Process) Run() {
	p.Cmd = exec.Command("java", "-Djline.terminal=jline.UnsupportedTerminal", "-jar", "BungeeCord.jar")
	//p.Cmd = exec.Command("bash", "fake-server.sh")
	p.Cmd.Dir = "bungee"

	stdout, _ := p.Cmd.StdoutPipe()
	stderr, _ := p.Cmd.StderrPipe()
	stdin, _ := p.Cmd.StdinPipe()

	if err := p.Cmd.Start(); err != nil {
		log.Fatal("cmd.Start()", err)
	}

	p.Running = true

	var wg sync.WaitGroup
	wg.Add(4)
	go func() {
		defer wg.Done()

		for {
			msg := <-p.Input
			io.WriteString(stdin, msg+"\n")
		}
	}()

	go func() {
		defer wg.Done()

		in := bufio.NewScanner(stdout)

		for in.Scan() {
			p.Output <- in.Text()
		}

		if err := in.Err(); err != nil {
			log.Printf("error: %s", err)
		}
	}()

	go func() {
		defer wg.Done()

		in2 := bufio.NewScanner(stderr)
		for in2.Scan() {
			log.Info(in2.Text())
		}
	}()

	go func() {
		defer wg.Done()

		if err := p.Cmd.Wait(); err != nil {
			if exiterr, ok := err.(*exec.ExitError); ok {
				if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
					/*
						Exit status: 143 == SIGTERM
						Exit status: -1  == SIGKILL?
					*/
					log.Printf("Exit status: %d", status.ExitStatus())
				}
			}
			log.Errorf("cmd.Wait: %v", err)
		} else {
			os.Exit(0)
		}
	}()

	wg.Wait()
}

func (p *Process) Log() {
	for {
		msg := <-p.Output

		logMessage := lib.LogMessage{
			Message: msg,
		}

		logMessageJson, _ := json.Marshal(logMessage)

		err := p.Channel.Publish("", p.QueueNormal.Name, false, false, amqp.Publishing{
			ContentType: "application/json",
			Body:        logMessageJson,
		})
		lib.FailOnError(err, "Publish error")

		log.WithFields(log.Fields{
			"msg": msg,
		}).Info()
	}
}

func (p *Process) Kill(signal syscall.Signal) {
	if p.Running {
		p.Cmd.Process.Signal(signal)
		p.Running = false
	}

}

func (p *Process) Rpc() {
	msgs, err := p.Channel.Consume(p.RpcQueue.Name, "", false, false, false, false, nil)
	lib.FailOnError(err, "Failed to register a consumer")

	for d := range msgs {
		log.Info("Got RPC call from RabbitMQ")
		var wr lib.RpcMessage
		err := json.Unmarshal(d.Body, &wr)
		if err != nil {
			log.Warn(err)
		}

		switch wr.Type {
		case lib.START:
			log.Info("START message")
			go p.Run()

			err = p.Channel.Publish("", d.ReplyTo, false, false, amqp.Publishing{
				ContentType:   "application/json",
				CorrelationId: d.CorrelationId,
				Body:          []byte(""),
			})
			lib.FailOnError(err, "Failed to publish a message")

			d.Ack(false)
		case lib.COMMAND:
			log.Info("COMMAND message")

			go func() { p.Input <- string(wr.Body) }()

			err = p.Channel.Publish("", d.ReplyTo, false, false, amqp.Publishing{
				ContentType:   "application/json",
				CorrelationId: d.CorrelationId,
				Body:          []byte(""),
			})
			lib.FailOnError(err, "Failed to publish a message")

			d.Ack(false)
		case lib.STOP_SIGKILL:
			log.Info("STOP_SIGKILL message")

			p.Kill(syscall.SIGKILL)

			err = p.Channel.Publish("", d.ReplyTo, false, false, amqp.Publishing{
				ContentType:   "application/json",
				CorrelationId: d.CorrelationId,
				Body:          []byte(""),
			})
			lib.FailOnError(err, "Failed to publish a message")

			d.Ack(false)
		case lib.STOP_SIGTERM:
			log.Info("STOP_SIGTERM message")

			p.Kill(syscall.SIGTERM)

			err = p.Channel.Publish("", d.ReplyTo, false, false, amqp.Publishing{
				ContentType:   "application/json",
				CorrelationId: d.CorrelationId,
				Body:          []byte(""),
			})
			lib.FailOnError(err, "Failed to publish a message")

			d.Ack(false)
		}
	}
}
