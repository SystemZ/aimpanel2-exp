package wrapper

import (
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/lib/task"
	"gitlab.com/systemz/aimpanel2/slave/config"
	"gitlab.com/systemz/aimpanel2/slave/model"
	"log"
	"os"
	"strings"
	"syscall"
)

func (p *Process) RedisListener(done chan bool) {
	// start connection to redis
	model.InitRedis()

	// subscribe tasks
	// https://godoc.org/github.com/go-redis/redis#PubSub
	pubsub := model.Redis.Subscribe(config.REDIS_PUB_SUB_WRAPPER_CH)

	// Wait for confirmation that subscription is created before publishing anything.
	_, err := pubsub.Receive()
	if err != nil {
		// FIXME don't panic on redis pub/sub error
		panic(err)
	}
	defer pubsub.Close()

	// Go channel which receives messages.
	ch := pubsub.Channel()

	done <- true

	// Consume messages.
	for msg := range ch {
		p.RedisTaskHandler(msg.Channel, msg.Payload)
	}

}

func (p *Process) RedisTaskHandler(taskCh string, taskBody string) {
	taskMsg := task.Message{}
	err := taskMsg.Deserialize(taskBody)
	if err != nil {
		logrus.Error(err)
	}

	// accept message only for our game servers or all on host
	if taskMsg.GameServerID != p.GameServerID && taskMsg.GameServerID != "all" {
		log.Printf("wrapper task is not for me, ignoring...")
		return
	}

	switch taskMsg.TaskId {
	case task.GAME_START:
		logrus.Infof("Got %v", taskMsg.TaskId)

		startCommand, err := taskMsg.Game.GetCmd()
		if err != nil {
			logrus.Error(err)
		}

		p.GameStartCommand = strings.Split(startCommand, " ")

		go p.Run()
	case task.GAME_STOP_SIGTERM:
		logrus.Infof("Got %v", taskMsg.TaskId)
		p.Kill(syscall.SIGTERM)
		os.Exit(0)
	case task.GAME_STOP_SIGKILL:
		logrus.Infof("Got %v", taskMsg.TaskId)
		p.Kill(syscall.SIGKILL)
		os.Exit(0)
	case task.GAME_COMMAND:
		logrus.Infof("Got %v", taskMsg.TaskId)
		go func() { p.Input <- taskMsg.Body }()
	case task.GAME_METRICS_FREQUENCY:
		logrus.Infof("Got %v", taskMsg.TaskId)
		p.MetricFrequency = taskMsg.MetricFrequency
		go p.Metrics()
	default:
		logrus.Warningf("Unhandled task %v!", taskMsg.TaskId)
		// report this to master

	}
}
