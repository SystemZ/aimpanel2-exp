package host

import (
	"github.com/alexandrevicenzi/go-sse"
	"github.com/dgrijalva/jwt-go"
	"gitlab.com/systemz/aimpanel2/lib"
	"gitlab.com/systemz/aimpanel2/lib/ecode"
	"gitlab.com/systemz/aimpanel2/lib/request"
	"gitlab.com/systemz/aimpanel2/lib/task"
	"gitlab.com/systemz/aimpanel2/master/events"
	"gitlab.com/systemz/aimpanel2/master/model"
	"gitlab.com/systemz/aimpanel2/master/service/gameserver"
	"os"
	"time"
)

func Create(data *request.HostCreate, userId string) (*model.Host, int) {
	host := &model.Host{
		Base: model.Base{
			DocType: "host",
		},
		Name:            data.Name,
		Ip:              data.Ip,
		UserId:          userId,
		MetricFrequency: 30,
		Token:           lib.RandomString(32),
	}
	err := host.Put(&host)
	if err != nil {
		return nil, ecode.DbSave
	}

	group := model.GetGroup("USER-" + userId)
	if group == nil {
		return nil, ecode.GroupNotFound
	}

	// FIXME handle errors
	model.CreatePermissionsForNewHost(group.ID, host.ID)

	return host, ecode.NoError
}

//Removes host and linked game servers
func Remove(hostId string) int {
	host := model.GetHost(hostId)
	gameServers := model.GetGameServersByHostId(hostId)
	for _, gameServer := range *gameServers {
		err := gameserver.Remove(gameServer.ID)
		if err != nil {
			return ecode.GsRemove
		}
	}

	permissions := model.GetPermisionsByEndpointRegex("/v1/host/" + host.ID)
	for _, perm := range permissions {
		err := model.Delete(perm.ID, perm.Rev)
		if err != nil {
			return ecode.DbError
		}
	}

	err := model.Delete(host.ID, host.Rev)
	if err != nil {
		return ecode.DbError
	}

	return ecode.NoError
}

func Auth(t string) (string, int) {
	host := model.GetHostByToken(t)

	if host == nil {
		return "", ecode.HostNotFound
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * 48).Unix(),
		"uid": host.ID,
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", ecode.Unknown
	}

	return tokenString, ecode.NoError
}

func CreateJob(data *request.HostCreateJob, userId string, hostId string) (*model.HostJob, int) {
	hostJob := &model.HostJob{
		Base: model.Base{
			DocType: "host_job",
		},
		Name:           data.Name,
		HostId:         hostId,
		CronExpression: data.CronExpression,
		TaskMessage:    data.TaskMessage,
	}

	err := hostJob.Put(&hostJob)
	if err != nil {
		return nil, ecode.DbSave
	}

	group := model.GetGroup("USER-" + userId)
	if group == nil {
		return nil, ecode.GroupNotFound
	}

	// FIXME handle errors
	model.CreatePermissionsForNewHostJob(group.ID, hostId, hostJob.ID)

	ec := sendJobsToAgent(hostId)
	if ec != ecode.NoError {
		return nil, ec
	}

	return hostJob, ecode.NoError
}

func RemoveJob(hostId string, jobId string) int {
	hostJob := model.GetHostJob(jobId)

	permissions := model.GetPermisionsByEndpointRegex("/v1/host/" + hostId + "/job/" + jobId)
	for _, perm := range permissions {
		err := model.Delete(perm.ID, perm.Rev)
		if err != nil {
			return ecode.DbError
		}
	}

	err := model.Delete(hostJob.ID, hostJob.Rev)
	if err != nil {
		return ecode.DbError
	}

	return sendJobsToAgent(hostId)
}

func sendJobsToAgent(hostId string) int {
	host := model.GetHost(hostId)
	if host == nil {
		return ecode.HostNotFound
	}

	var jobs []task.Job

	hostJobs := model.GetHostJobs(host.ID)
	for _, job := range hostJobs {
		jobs = append(jobs, task.Job{
			Name:           job.Name,
			CronExpression: job.CronExpression,
			TaskMessage:    job.TaskMessage,
		})
	}

	channel, ok := events.SSE.GetChannel("/v1/events/" + host.Token)
	if !ok {
		return ecode.HostNotTurnedOn
	}

	taskMsg := task.Message{
		TaskId: task.AGENT_GET_JOBS,
		Jobs:   &jobs,
	}

	taskMsgStr, err := taskMsg.Serialize()
	if err != nil {
		return ecode.Unknown
	}

	channel.SendMessage(sse.NewMessage("", taskMsgStr, taskMsg.TaskId.StringValue()))

	return ecode.NoError
}
