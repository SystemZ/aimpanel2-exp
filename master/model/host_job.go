package model

type HostJob struct {
	Base

	//User assigned name
	Name string `json:"name" example:"Restart server"`

	// Host ID
	//
	// required: true
	HostId string `json:"host_id" example:"100112233-4455-6677-8899-aabbccddeeff"`

	// Game server ID
	//
	// required: false
	GameServerId string `json:"game_server_id"`

	CronExpression string `json:"cron_expression" example:"5 4 * * *"`

	TaskMessageJson string `json:"task_message_json"`
}

func GetHostJobs(hostId string) []HostJob {
	var hj []HostJob

	err := GetS(&hj, map[string]interface{}{
		"doc_type": "host_job",
	})
	if err != nil {
		return nil
	}

	return hj
}
