package agent

import (
	"fmt"
	"github.com/coreos/go-systemd/sdjournal"
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/lib"
	"gitlab.com/systemz/aimpanel2/lib/ahttp"
	"gitlab.com/systemz/aimpanel2/lib/response"
	"gitlab.com/systemz/aimpanel2/lib/task"
	"gitlab.com/systemz/aimpanel2/slave/config"
	"gitlab.com/systemz/aimpanel2/slave/cron"
	"gitlab.com/systemz/aimpanel2/slave/model"
	"gitlab.com/systemz/aimpanel2/slave/tasks"
	"math"
	"os"
	"time"
)

func Start(hostToken string) {
	model.InitRedis()
	cron.InitCron()
	ahttp.HttpClient = ahttp.InitHttpClient()

	logrus.Info("Starting Agent Version." + config.GIT_COMMIT)
	config.HOST_TOKEN = hostToken

	var token response.Token
	_, err := ahttp.Get(config.API_URL+"/v1/host/auth/"+config.HOST_TOKEN, &token)
	if err != nil {
		lib.FailOnError(err, "Failed to get host token")
	}
	config.API_TOKEN = token.Token

	sseStarted := make(chan bool, 1)
	redisStarted := make(chan bool, 1)
	// all tasks from master are handled here
	go listenerSse(sseStarted)
	// wrapper and cli handling
	go listenerRedis(redisStarted)

	<-sseStarted
	<-redisStarted

	logrus.Info("Send AGENT_STARTED")
	taskMsg := task.Message{
		TaskId: task.AGENT_STARTED,
	}
	//TODO: do something with status code
	_, err = ahttp.SendTaskData(config.API_URL+"/v1/events/"+config.HOST_TOKEN, config.API_TOKEN, taskMsg)
	if err != nil {
		logrus.Error(err)
	}

	logrus.Info("Send AGENT_METRICS_FREQUENCY")
	taskMsg = task.Message{
		TaskId: task.AGENT_METRICS_FREQUENCY,
	}
	//TODO: do something with status code
	_, err = ahttp.SendTaskData(config.API_URL+"/v1/events/"+config.HOST_TOKEN, config.API_TOKEN, taskMsg)
	if err != nil {
		logrus.Error(err)
	}

	go func() {
		time.Sleep(time.Duration(lib.RandInt(200, 2000)) * time.Millisecond)

		logrus.Info("Send AGENT_GET_JOBS")
		taskMsg = task.Message{
			TaskId: task.AGENT_GET_JOBS,
		}

		//TODO: do something with status code
		_, err = ahttp.SendTaskData(config.API_URL+"/v1/events/"+config.HOST_TOKEN, config.API_TOKEN, taskMsg)
		if err != nil {
			logrus.Error(err)
		}
	}()

	go func() {
		logrus.Info("Starting ssh log parsing...")
		//journalInstance, err := sdjournal.NewJournal()
		jr, err := sdjournal.NewJournalReader(sdjournal.JournalReaderConfig{
			Since: -time.Duration(5) * time.Second,
			//Since: 1 * time.Nanosecond,
			Formatter: func(entry *sdjournal.JournalEntry) (string, error) {
				logrus.WithFields(logrus.Fields{
					"entry":   *entry,
					"message": entry.Fields["MESSAGE"],
				}).Debug("Message from journal received")
				return fmt.Sprintln(entry.Fields["MESSAGE"]), nil
			},
			//Path:  "/var/log/journal",
			//Path:  "/run/log/journal", // not for ubuntu
			//NumFromTail: 10,
			//Cursor:
			Matches: []sdjournal.Match{
				{
					Field: sdjournal.SD_JOURNAL_FIELD_SYSTEMD_UNIT,
					Value: "ssh.service",
				},
			},
			//Matches: []sdjournal.Match{
			//	{"SYSLOG_IDENTIFIER", "ssh.service"},
			//},
		})

		if err != nil {
			logrus.Error(err)
		}

		if jr == nil {
			logrus.Error("nil journal reader")
		}

		defer jr.Close()
		err = jr.Follow(time.After(time.Duration(math.MaxInt64)), os.Stdout)
		if err != nil {
			logrus.Error(err)
		}

		/*
			b := make([]byte, 64*1<<(10)) // 64KB.
			for {
				c, err := jr.Read(b)
				if err != nil {
					if err == io.EOF {
						break
					}
					panic(err)
				}
				logrus.Info(string(b[:c]))
			}
		*/

		logrus.Info("Stopped ssh log parsing...")
	}()

	tasks.AgentSendOSInfo()

	select {}
}
