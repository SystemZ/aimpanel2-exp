package cron

import (
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/lib/task"
	"gitlab.com/systemz/aimpanel2/slave/tasks"
)

var Cron *cron.Cron

func InitCron() {
	Cron = cron.New()
	logrus.Info("Cron initialized")
}

func AddJobs(jobs []task.Job) {
	Stop()

	for _, job := range jobs {
		logrus.Infof("Adding job %v", job.Name)
		_, err := Cron.AddFunc(job.CronExpression, func() {
			tasks.ProcessTask(job.TaskMessage)
		})
		if err != nil {
			logrus.Error("Error while setting up %v cron job", job.Name)
		}
	}

	Cron.Start()
}

func Stop() {
	Cron.Stop()
}
