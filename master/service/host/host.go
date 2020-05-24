package host

import (
	"errors"
	"github.com/alexandrevicenzi/go-sse"
	"github.com/dgrijalva/jwt-go"
	"gitlab.com/systemz/aimpanel2/lib"
	"gitlab.com/systemz/aimpanel2/lib/ecode"
	"gitlab.com/systemz/aimpanel2/lib/request"
	"gitlab.com/systemz/aimpanel2/lib/task"
	"gitlab.com/systemz/aimpanel2/master/events"
	"gitlab.com/systemz/aimpanel2/master/model"
	"gitlab.com/systemz/aimpanel2/master/service/gameserver"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
	"time"
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

	group, err := model.GetGroupByName("USER-" + userId.Hex())
	if err != nil {
		return nil, ecode.DbError
	}

	if group == nil {
		return nil, ecode.GroupNotFound
	}

	// FIXME handle errors
	model.CreatePermissionsForNewHost(group.ID, host.ID)

	return host, ecode.NoError
}

//Removes host and linked game servers
func Remove(hostId primitive.ObjectID) int {
	host, err := model.GetHostById(hostId)
	if err != nil {
		return ecode.DbError
	}

	gameServers, err := model.GetGameServersByHostId(hostId)
	if err != nil {
		return ecode.DbError
	}

	for _, gameServer := range *gameServers {
		err := gameserver.Remove(gameServer.ID)
		if err != nil {
			return ecode.GsRemove
		}
	}

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

	err = model.Delete(host)
	if err != nil {
		return ecode.DbError
	}

	return ecode.NoError
}

func Auth(t string) (string, int) {
	host, err := model.GetHostByToken(t)
	if err != nil {
		return "", ecode.DbError
	}

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

	group, err := model.GetGroupByName("USER-" + userId.Hex())
	if err != nil {
		return nil, ecode.DbError
	}

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

func Update(hostId primitive.ObjectID) error {
	hostToken, err := model.GetHostTokenById(hostId)
	if err != nil {
		return err
	}

	if hostToken == "" {
		return errors.New("error when getting host token from db")
	}

	channel, ok := events.SSE.GetChannel("/v1/events/" + hostToken)
	if !ok {
		return errors.New("host is not turned on")
	}

	commit, err := model.GetSlaveCommit(model.Redis)
	if err != nil {
		return err
	}

	url, err := model.GetSlaveUrl(model.Redis)
	if err != nil {
		return err
	}

	taskMsg := task.Message{
		TaskId: task.AGENT_UPDATE,
		Commit: commit,
		Url:    url,
	}

	taskMsgStr, err := taskMsg.Serialize()
	if err != nil {
		return err
	}

	channel.SendMessage(sse.NewMessage("", taskMsgStr, taskMsg.TaskId.StringValue()))

	return nil
}
