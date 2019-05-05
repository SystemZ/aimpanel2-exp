package wrapper

import (
	"bufio"
	"encoding/json"
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"gitlab.com/systemz/aimpanel2/lib"
	"gitlab.com/systemz/aimpanel2/lib/rabbit"
	"gitlab.com/systemz/aimpanel2/slave/config"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
)

type Process struct {
	Cmd     *exec.Cmd
	Running bool

	Output chan string
	Input  chan string

	//amqp
	Channel             *amqp.Channel
	Queue               amqp.Queue
	QueueLogs           amqp.Queue
	ClientCorrelationId string
	ReplyTo             string

	//
	GameServerID     string
	GameStartCommand []string
}

func (p *Process) Run() {
	p.Cmd = exec.Command(p.GameStartCommand[0], p.GameStartCommand[1:]...)
	p.Cmd.Dir = config.GS_DIR + p.GameServerID

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
			p.LogStdout(in.Text())
			logrus.Info(in.Text())
		}
	}()

	go func() {
		defer wg.Done()

		in := bufio.NewScanner(stderr)

		for in.Scan() {
			p.LogStderr(in.Text())
			logrus.Info(in.Text())
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
					//Wrapper think that this is a GAME START command
					//exitMessage := rabbit.ExitMessage{
					//	Code:    status.ExitStatus(),
					//	Message: "",
					//}
					//exitMessageJson, _ := json.Marshal(exitMessage)
					//
					//err := p.Channel.Publish("", p.Queue.Name, false, false, amqp.Publishing{
					//	ContentType: "application/json",
					//	Body:        exitMessageJson,
					//})
					//lib.FailOnError(err, "Publish error")

					log.Printf("Exit status: %d", status.ExitStatus())
				}
			}
			logrus.Errorf("cmd.Wait: %v", err)
			p.WrapperExitMessage()
			os.Exit(0)
		} else {
			p.WrapperExitMessage()
			os.Exit(0)
		}
	}()

	wg.Wait()
	logrus.Info("WG Done")
}

func (p *Process) LogStdout(msg string) {
	logMessage := rabbit.QueueMsg{
		TaskId: rabbit.SERVER_LOG,
		Stdout:       msg,
		GameServerID: uuid.FromStringOrNil(p.GameServerID),
	}

	logMessageJson, _ := json.Marshal(logMessage)

	err := p.Channel.Publish(
		"",
		p.QueueLogs.Name,
		false,
		false,
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: p.ClientCorrelationId,
			Body:          logMessageJson,
		})

	lib.FailOnError(err, "Publish error")
}

func (p *Process) LogStderr(msg string) {
	logMessage := rabbit.QueueMsg{
		TaskId: rabbit.SERVER_LOG,
		Stderr:       msg,
		GameServerID: uuid.FromStringOrNil(p.GameServerID),
	}

	logMessageJson, _ := json.Marshal(logMessage)

	logrus.Info(logMessage)

	err := p.Channel.Publish(
		"",
		p.QueueLogs.Name,
		false,
		false,
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: p.ClientCorrelationId,
			Body:          logMessageJson,
		})

	lib.FailOnError(err, "Publish error")
}

func (p *Process) Kill(signal syscall.Signal) {
	logrus.Info("Kill" + signal.String())
	if p.Running {
		p.Cmd.Process.Signal(signal)
		p.Running = false
	}

}

type rabbitTask struct {
	msg     amqp.Delivery
	msgBody rabbit.QueueMsg
	ch      *amqp.Channel
}

func (p *Process) Rpc() {
	msgs, err := p.Channel.Consume(p.Queue.Name, "", false, false, false, false, nil)
	lib.FailOnError(err, "Failed to register a consumer")

	for msg := range msgs {
		logrus.Info("Received a task")

		var msgBody rabbit.QueueMsg
		err := json.Unmarshal(msg.Body, &msgBody)
		if err != nil {
			logrus.Warn(err)
		}

		task := rabbitTask{
			msg:     msg,
			ch:      p.Channel,
			msgBody: msgBody,
		}

		switch msgBody.TaskId {
		case rabbit.GAME_START:
			logrus.Info("Got GAME_START msg")

			p.GameStartCommand = strings.Split(task.msgBody.GameStartCommand.Command, " ")

			go p.Run()

			err = p.Channel.Publish("", msg.ReplyTo, false, false, amqp.Publishing{
				ContentType:   "application/json",
				CorrelationId: p.ClientCorrelationId,
				Body:          []byte(""),
			})
			lib.FailOnError(err, "Failed to publish a message")

			msg.Ack(false)
		case rabbit.GAME_COMMAND:
			logrus.Info("Got GAME_COMMAND msg")

			go func() { p.Input <- string(msgBody.Body) }()

			err = p.Channel.Publish("", msg.ReplyTo, false, false, amqp.Publishing{
				ContentType:   "application/json",
				CorrelationId: p.ClientCorrelationId,
				Body:          []byte(""),
			})
			lib.FailOnError(err, "Failed to publish a message")

			msg.Ack(false)
		case rabbit.GAME_STOP_SIGKILL:
			logrus.Info("Got GAME_STOP_SIGKILL msg")

			p.Kill(syscall.SIGKILL)

			err = p.Channel.Publish("", msg.ReplyTo, false, false, amqp.Publishing{
				ContentType:   "application/json",
				CorrelationId: p.ClientCorrelationId,
				Body:          []byte(""),
			})
			lib.FailOnError(err, "Failed to publish a message")

			msg.Ack(false)
		case rabbit.GAME_STOP_SIGTERM:
			logrus.Info("Got GAME_STOP_SIGTERM msg")

			p.Kill(syscall.SIGTERM)

			err = p.Channel.Publish("", msg.ReplyTo, false, false, amqp.Publishing{
				ContentType:   "application/json",
				CorrelationId: p.ClientCorrelationId,
				Body:          []byte(""),
			})
			lib.FailOnError(err, "Failed to publish a message")

			msg.Ack(false)
		}
	}
}


func (p *Process) WrapperExitMessage() {
	msg := rabbit.QueueMsg{
		TaskId: rabbit.WRAPPER_EXITED,
		GameServerID: uuid.FromStringOrNil(p.GameServerID),
	}

	msgJson, _ := json.Marshal(msg)

	err := p.Channel.Publish(
		"",
		p.QueueLogs.Name,
		false,
		false,
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: p.ClientCorrelationId,
			Body:          msgJson,
		})

	lib.FailOnError(err, "Publish error")
}

func (p *Process) WrapperStartMessage() {
	logrus.Info("Sending WRAPPER_STARTED")
	msg := rabbit.QueueMsg{
		TaskId: rabbit.WRAPPER_STARTED,
		GameServerID: uuid.FromStringOrNil(p.GameServerID),
	}

	msgJson, _ := json.Marshal(msg)

	err := p.Channel.Publish(
		"",
		p.QueueLogs.Name,
		false,
		false,
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: p.ClientCorrelationId,
			Body:          msgJson,
		})

	lib.FailOnError(err, "Publish error")
}