package migrations

import (
	"gitlab.com/systemz/aimpanel2/lib/task"
	"gitlab.com/systemz/aimpanel2/master/model"
)

func Migration4Up() (err error) {
	hosts, err := model.GetHosts()
	if err != nil {
		return
	}

	for _, host := range hosts {
		hostJob := &model.HostJob{
			Name:           "Check for update",
			HostId:         host.ID,
			CronExpression: "*/10 * * * *",
			TaskMessage: task.Message{
				TaskId: task.AGENT_GET_UPDATE,
			},
		}
		err = model.Put(hostJob)
		if err != nil {
			return
		}
	}
	return
}

func Migration4Down() {

}
