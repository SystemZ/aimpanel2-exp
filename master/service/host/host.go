package host

import (
	"gitlab.com/systemz/aimpanel2/lib"
	"gitlab.com/systemz/aimpanel2/lib/ecode"
	"gitlab.com/systemz/aimpanel2/lib/request"
	"gitlab.com/systemz/aimpanel2/lib/task"
	"gitlab.com/systemz/aimpanel2/master/config"
	"gitlab.com/systemz/aimpanel2/master/model"
	"gitlab.com/systemz/aimpanel2/master/service/gameserver"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Create(data *request.HostCreate, userId primitive.ObjectID) (*model.Host, int) {
	host := &model.Host{
		Name:            data.Name,
		Ip:              data.Ip,
		UserId:          userId,
		MetricFrequency: 30,
		Token:           lib.RandomString(32),
	}
	err := model.Put(host)
	if err != nil {
		return nil, ecode.DbSave
	}

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
		return nil, ecode.DbSave
	}

	/*
		err = model.CreatePermissionsForNewHost(userId, host.ID)
		if err != nil {
			return nil, ecode.DbSave
		}
	*/

	return host, ecode.NoError
}

//Removes host and linked game servers
func Remove(hostId primitive.ObjectID, user model.User) int {
	host, err := model.GetHostById(hostId)
	if err != nil {
		return ecode.DbError
	}

	gameServers, err := model.GetGameServersByHostId(hostId)
	if err != nil {
		return ecode.DbError
	}

	for _, gameServer := range *gameServers {
		err := gameserver.Remove(gameServer.ID, user)
		if err != nil {
			return ecode.GsRemove
		}
	}

	/*
		permissions, err := model.GetPermisionsByEndpointRegex("/v1/host/" + host.ID.Hex())
		if err != nil {
			return ecode.DbError
		}

		for _, perm := range permissions {
			err := model.Delete(&perm)
			if err != nil {
				return ecode.DbError
			}
		}
	*/

	err = model.Delete(host)
	if err != nil {
		return ecode.DbError
	}

	return ecode.NoError
}

func CreateJob(data *request.HostCreateJob, userId primitive.ObjectID, hostId primitive.ObjectID) (*model.HostJob, int) {
	hostJob := &model.HostJob{
		Name:           data.Name,
		HostId:         hostId,
		CronExpression: data.CronExpression,
		TaskMessage:    data.TaskMessage,
	}

	err := model.Put(hostJob)
	if err != nil {
		return nil, ecode.DbSave
	}

	/*
		err = model.CreatePermissionsForNewHostJob(userId, hostId, hostJob.ID)
		if err != nil {
			return nil, ecode.DbSave
		}
	*/

	ec := sendJobsToAgent(hostId)
	if ec != ecode.NoError {
		return nil, ec
	}

	return hostJob, ecode.NoError
}

func RemoveJob(hostId primitive.ObjectID, jobId primitive.ObjectID) int {
	hostJob, err := model.GetHostJobById(jobId)
	if err != nil {
		return ecode.DbError
	}

	permissions, err := model.GetPermisionsByEndpointRegex("/v1/host/" + hostId.Hex() + "/job/" + jobId.Hex())
	if err != nil {
		return ecode.DbError
	}

	for _, perm := range permissions {
		err := model.Delete(&perm)
		if err != nil {
			return ecode.DbError
		}
	}

	err = model.Delete(hostJob)
	if err != nil {
		return ecode.DbError
	}

	return sendJobsToAgent(hostId)
}

func sendJobsToAgent(hostId primitive.ObjectID) int {
	host, err := model.GetHostById(hostId)
	if err != nil {
		return ecode.DbError
	}

	if host == nil {
		return ecode.HostNotFound
	}

	var jobs []task.Job

	hostJobs, err := model.GetHostJobsByHostId(host.ID)
	if err != nil {
		return ecode.DbError
	}

	for _, job := range hostJobs {
		jobs = append(jobs, task.Job{
			Name:           job.Name,
			CronExpression: job.CronExpression,
			TaskMessage:    job.TaskMessage,
		})
	}

	taskMsg := task.Message{
		TaskId: task.AGENT_GET_JOBS,
		Jobs:   &jobs,
	}

	err = model.SendTaskToSlave(host.ID, model.User{}, taskMsg)
	if err != nil {
		return ecode.DbSave
	}

	return ecode.NoError
}

func Update(hostId primitive.ObjectID, user model.User) error {
	host, err := model.GetHostById(hostId)
	if err != nil {
		return &lib.Error{ErrorCode: ecode.HostNotFound}
	}

	binaryUrl := config.HTTP_REPO_URL + "/aimpanel/latest"
	taskMsg := task.Message{
		TaskId: task.AGENT_UPDATE,
		// TODO in DEV mode just use some random string
		Commit: "new", // FIXME set git commit ID in master when building in CI
		Url:    binaryUrl,
	}

	err = model.SendTaskToSlave(host.ID, user, taskMsg)
	if err != nil {
		return &lib.Error{ErrorCode: ecode.DbSave}
	}

	return nil
}
