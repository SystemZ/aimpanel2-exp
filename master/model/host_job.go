package model

import (
	"gitlab.com/systemz/aimpanel2/lib/task"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HostJob struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty" example:"1238206236281802752"`

	//User assigned name
	Name string `json:"name" example:"Restart server"`

	// Host ID
	//
	// required: true
	HostId primitive.ObjectID `json:"host_id" example:"100112233-4455-6677-8899-aabbccddeeff"`

	CronExpression string `json:"cron_expression" example:"5 4 * * *"`

	TaskMessage task.Message `json:"task_message"`
}

func (h *HostJob) GetCollectionName() string {
	return "hosts_job"
}

func GetHostJobs(hostId primitive.ObjectID) []HostJob {
	var hj []HostJob

	err := GetS(&hj, map[string]interface{}{
		"doc_type": "host_job",
		"host_id":  hostId,
	})
	if err != nil {
		return nil
	}

	return hj
}

func GetHostJob(jobId primitive.ObjectID) *HostJob {
	var hostJob HostJob
	err := GetOneS(&hostJob, map[string]interface{}{
		"doc_type": "host_job",
		"_id":      jobId,
	})
	if err != nil {
		return nil
	}

	return &hostJob
}
