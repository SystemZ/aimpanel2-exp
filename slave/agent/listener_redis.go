package agent

import (
	"context"
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/lib/ahttp"
	"gitlab.com/systemz/aimpanel2/lib/task"
	"gitlab.com/systemz/aimpanel2/slave/config"
	"gitlab.com/systemz/aimpanel2/slave/cron"
	"gitlab.com/systemz/aimpanel2/slave/model"
	"gitlab.com/systemz/aimpanel2/slave/tasks"
	"sync"
	"time"
)

var (
	QueuedMsgs []task.Message
	QueueMutex sync.Mutex
)

func QueueSendTaskData(msgRaw task.Message) {
	QueueMutex.Lock()
	QueuedMsgs = append(QueuedMsgs, msgRaw)
	QueueMutex.Unlock()
}

// gather all messages in specified time and send them in batches
// prevent massive number of HTTP requests which wrongfully can look like DoS
// TODO put some kind of debouncer
// TODO use less sleep and check for high priority messages
// TODO decide sleep time with last sent nanosecond timestamp
func SendMessagesToMaster() {
	for {
		// wait between sending batches
		time.Sleep(time.Millisecond * 300)
		// lock for consistency
		QueueMutex.Lock()
		// no messages, no need to send them
		if len(QueuedMsgs) < 1 {
			QueueMutex.Unlock()
			continue
		}

		// send all messages in queue
		// FIXME handle task send retry
		_, err := ahttp.SendTaskBatchData("/v1/events/"+config.HOST_TOKEN+"/batch", config.HW_ID, QueuedMsgs)
		if err != nil {
			logrus.Error(err)
		}

		// clean queue
		QueuedMsgs = []task.Message{}
		// allow to add new messages for a second
		QueueMutex.Unlock()
	}
}

func listenerRedis(done chan bool) {
	// start connection to redis
	//model.InitRedis()

	// start batch processing
	go SendMessagesToMaster()

	// subscribe tasks
	// https://godoc.org/github.com/go-redis/redis#PubSub
	pubsub := model.Redis.Subscribe(context.TODO(), config.REDIS_PUB_SUB_AGENT_CH)

	// Wait for confirmation that subscription is created before publishing anything.
	_, err := pubsub.Receive(context.TODO())
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
		redisTaskHandler(msg.Channel, msg.Payload)
	}
}

func redisTaskHandler(taskCh string, taskBody string) {
	taskMsg := task.Message{}
	err := taskMsg.Deserialize(taskBody)
	if err != nil {
		logrus.Error(err)
	}

	switch taskMsg.TaskId {
	case task.GAME_STARTED:
		logrus.Info("Agent got " + taskMsg.TaskId.String())

		val, err := model.GetGsStart(taskMsg.GameServerID)
		if val == 1 {
			model.SetGsStart(taskMsg.GameServerID, 2)
			tasks.GsStartGame(taskMsg)
			model.DelGsStart(taskMsg.GameServerID)

			_, err = ahttp.SendTaskData("/v1/events/"+config.HOST_TOKEN, config.HW_ID, taskMsg)
			if err != nil {
				logrus.Error(err)
			}
		}

	case task.GAME_SHUTDOWN:
		logrus.Info("Agent got " + taskMsg.TaskId.String())

		_, err = ahttp.SendTaskData("/v1/events/"+config.HOST_TOKEN, config.HW_ID, taskMsg)
		if err != nil {
			logrus.Error(err)
		}

		//Start wrapper if gs is restarting
		// FIXME redesign restart cmd
		/*
			val, _ := model.GetGsRestart(taskMsg.GameServerID)
			if val == 1 {
				model.SetGsRestart(taskMsg.GameServerID, 2)
				tasks.StartWrapperInDocker(taskMsg.GameServerID)
				if false {
					tasks.StartWrapperExecRaw(taskMsg)
				}
				model.DelGsRestart(taskMsg.GameServerID)
			}
		*/

	case task.GAME_SERVER_LOG:
		logrus.Info("Agent got " + taskMsg.TaskId.String())
		QueueSendTaskData(taskMsg)
	case task.GAME_METRICS:
		logrus.Info("Agent got " + taskMsg.TaskId.String())

		_, err = ahttp.SendTaskData("/v1/events/"+config.HOST_TOKEN, config.HW_ID, taskMsg)
		if err != nil {
			logrus.Error(err)
		}
	case task.GAME_METRICS_FREQUENCY:
		logrus.Info("Agent got " + taskMsg.TaskId.String())

		_, err = ahttp.SendTaskData("/v1/events/"+config.HOST_TOKEN, config.HW_ID, taskMsg)
		if err != nil {
			logrus.Error(err)
		}

	case task.AGENT_BACKUP_GS:
		logrus.Info("agent got GAME_MAKE_BACKUP msg")
		go tasks.GsBackup(taskMsg.GameServerID)
	case task.AGENT_BACKUP_RESTORE_GS:
		logrus.Info("agent got GAME_BACKUP_RESTORE msg")
		go tasks.GsBackupRestore(taskMsg.GameServerID, taskMsg.BackupFilename)
	case task.AGENT_CLEAN_REINSTALL_GS:
		logrus.Info("agent got GAME_CLEAN_REINSTALL msg")
		go tasks.GsCleanReinstall(taskMsg)
	case task.AGENT_SHUTDOWN:
		logrus.Info("agent got AGENT_SHUTDOWN msg")
		cron.Stop()
		tasks.AgentShutdown()
	case task.SUPERVISOR_CLEAN_FILES_GS:
		model.GsCleanFilesFinished(taskMsg.GameServerID, true)
	default:
		logrus.Warningf("Unhandled task %v!", taskMsg.TaskId)
	}
}
