package wrapper

import (
	"bufio"
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/lib/task"
	"gitlab.com/systemz/aimpanel2/slave/config"
	"gitlab.com/systemz/aimpanel2/slave/model"
	"io"
	"log"
	"os"
	"os/exec"
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
					// FIXME verify launched app exit code and report it to agent
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
				TaskId: task.GAME_SHUTDOWN,
			}
			model.SendTask(config.REDIS_PUB_SUB_AGENT_CH, taskMsg)

			os.Exit(0)
		} else {
			os.Exit(0)
		}
	}()

	wg.Wait()
	logrus.Info("WG Done")
}

func (p *Process) Kill(signal syscall.Signal) {
	logrus.Info("Kill " + signal.String())
	if p.Running {
		p.Cmd.Process.Signal(signal)
		p.Running = false
	}

}
