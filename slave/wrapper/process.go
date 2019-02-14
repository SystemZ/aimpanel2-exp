package wrapper

import (
	"bufio"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"gitlab.com/systemz/aimpanel2/lib"
	"io"
	"log"
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

	//
	GameServerUUID string
	Game           lib.Game
}

func (p *Process) Run() {
	p.Cmd = exec.Command(p.Game.Command[0], p.Game.Command[1:]...)
	p.Cmd.Dir = "/opt/aimpanel/gs/" + p.GameServerUUID

	stdout, _ := p.Cmd.StdoutPipe()
	stderr, _ := p.Cmd.StderrPipe()
	stdin, _ := p.Cmd.StdinPipe()

	if err := p.Cmd.Start(); err != nil {
		logrus.Fatal("cmd.Start()", err)
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
			logrus.Printf("error: %s", err)
		}
	}()

	go func() {
		defer wg.Done()

		in2 := bufio.NewScanner(stderr)
		for in2.Scan() {
			logrus.Info(in2.Text())
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
					exitMessage := lib.ExitMessage{
						Code:    status.ExitStatus(),
						Message: "",
					}

					exitMessageJson, _ := json.Marshal(exitMessage)

					err := p.Channel.Publish("", p.QueueHigh.Name, false, false, amqp.Publishing{
						ContentType: "application/json",
						Body:        exitMessageJson,
					})
					lib.FailOnError(err, "Publish error")

					log.Printf("Exit status: %d", status.ExitStatus())
				}
			}
			logrus.Errorf("cmd.Wait: %v", err)

			os.Exit(0)
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

		logrus.WithFields(logrus.Fields{
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
		logrus.Info("Got RPC call from RabbitMQ")

		var rpcMsg lib.RpcMessage
		err := json.Unmarshal(d.Body, &rpcMsg)
		if err != nil {
			logrus.Warn(err)
		}

		switch rpcMsg.Type {
		case lib.GAME_START:
			logrus.Info("Got GAME_START msg")

			p.Game = lib.GAMES[rpcMsg.Game]
			p.GameServerUUID = rpcMsg.GameServerUUID

			go p.Run()

			err = p.Channel.Publish("", d.ReplyTo, false, false, amqp.Publishing{
				ContentType:   "application/json",
				CorrelationId: d.CorrelationId,
				Body:          []byte(""),
			})
			lib.FailOnError(err, "Failed to publish a message")

			d.Ack(false)
		case lib.GAME_COMMAND:
			logrus.Info("Got GAME_COMMAND msg")

			go func() { p.Input <- string(rpcMsg.Body) }()

			err = p.Channel.Publish("", d.ReplyTo, false, false, amqp.Publishing{
				ContentType:   "application/json",
				CorrelationId: d.CorrelationId,
				Body:          []byte(""),
			})
			lib.FailOnError(err, "Failed to publish a message")

			d.Ack(false)
		case lib.GAME_STOP_SIGKILL:
			logrus.Info("Got GAME_STOP_SIGKILL msg")

			p.Kill(syscall.SIGKILL)

			err = p.Channel.Publish("", d.ReplyTo, false, false, amqp.Publishing{
				ContentType:   "application/json",
				CorrelationId: d.CorrelationId,
				Body:          []byte(""),
			})
			lib.FailOnError(err, "Failed to publish a message")

			d.Ack(false)
		case lib.GAME_STOP_SIGTERM:
			logrus.Info("Got GAME_STOP_SIGTERM msg")

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
