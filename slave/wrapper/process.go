package wrapper

import (
	"bufio"
	"github.com/r3labs/sse"
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/lib"
	"gitlab.com/systemz/aimpanel2/lib/task"
	"gitlab.com/systemz/aimpanel2/slave/config"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"syscall"
)

type Process struct {
	Cmd     *exec.Cmd
	Running bool

	Output chan string
	Input  chan string

	GameServerID     string
	GameStartCommand []string
	MetricFrequency  int
}

func (p *Process) Run() {
	p.Cmd = exec.Command(p.GameStartCommand[0], p.GameStartCommand[1:]...)
	p.Cmd.Dir = config.GS_DIR + "/" + p.GameServerID

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

			taskMsg := task.Message{
				TaskId: task.WRAPPER_EXITED,
			}

			jsonStr, err := taskMsg.Serialize()
			if err != nil {
				logrus.Error(err)
			}
			//TODO: do something with status code
			_, err = lib.SendTaskData(config.API_URL+"/v1/events/"+config.HOST_TOKEN+"/"+p.GameServerID, config.API_TOKEN, jsonStr)
			if err != nil {
				logrus.Error(err)
			}
			os.Exit(0)
		} else {
			taskMsg := task.Message{
				TaskId: task.WRAPPER_EXITED,
			}

			jsonStr, err := taskMsg.Serialize()
			if err != nil {
				logrus.Error(err)
			}
			//TODO: do something with status code
			_, err = lib.SendTaskData(config.API_URL+"/v1/events/"+config.HOST_TOKEN+"/"+p.GameServerID, config.API_TOKEN, jsonStr)
			if err != nil {
				logrus.Error(err)
			}
			os.Exit(0)
		}
	}()

	wg.Wait()
	logrus.Info("WG Done")
}

func (p *Process) Kill(signal syscall.Signal) {
	logrus.Info("Kill" + signal.String())
	if p.Running {
		p.Cmd.Process.Signal(signal)
		p.Running = false
	}

}

func (p *Process) Rpc() {
	client := sse.NewClient(config.API_URL + "/v1/events/" + config.HOST_TOKEN + "/" + p.GameServerID)
	client.Headers = map[string]string{
		"Authorization": "Bearer " + config.API_TOKEN,
	}
	err := client.SubscribeRaw(func(msg *sse.Event) {
		logrus.Info(msg.ID)
		logrus.Info(string(msg.Data))
		logrus.Info(string(msg.Event))

		taskMsg := task.Message{}
		err := taskMsg.Deserialize(string(msg.Data))
		if err != nil {
			logrus.Error(err)
		}

		taskId, _ := strconv.Atoi(string(msg.Event))

		switch taskId {
		case task.GAME_START:
			logrus.Info("Got GAME_START msg")

			startCommand, err := taskMsg.Game.GetCmd()
			if err != nil {
				logrus.Error(err)
			}

			p.GameStartCommand = strings.Split(startCommand, " ")

			go p.Run()
		case task.GAME_COMMAND:
			logrus.Info("Got GAME_COMMAND msg")
			go func() { p.Input <- taskMsg.Body }()
		case task.GAME_STOP_SIGKILL:
			logrus.Info("Got GAME_STOP_SIGKILL msg")
			p.Kill(syscall.SIGKILL)
		case task.GAME_STOP_SIGTERM:
			logrus.Info("Got GAME_STOP_SIGTERM msg")
			p.Kill(syscall.SIGTERM)
		case task.WRAPPER_METRICS_FREQUENCY:
			logrus.Info("Got WRAPPER_METRICS_FREQUENCY msg")
			p.MetricFrequency = taskMsg.MetricFrequency
			go p.Metrics()
		}
	})
	if err != nil {
		lib.FailOnError(err, "Can't connect to event channel")
	}
}
