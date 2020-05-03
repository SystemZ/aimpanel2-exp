package gameserver

import (
	"encoding/json"
	"errors"
	"github.com/alexandrevicenzi/go-sse"
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/lib"
	"gitlab.com/systemz/aimpanel2/lib/ecode"
	"gitlab.com/systemz/aimpanel2/lib/filemanager"
	"gitlab.com/systemz/aimpanel2/lib/game"
	"gitlab.com/systemz/aimpanel2/lib/task"
	"gitlab.com/systemz/aimpanel2/master/events"
	"gitlab.com/systemz/aimpanel2/master/model"
	"time"
)

func Start(gsId string) error {
	gameServer := model.GetGameServer(gsId)
	if gameServer == nil {
		return errors.New("getting game server from db failed")
	}

	hostToken := model.GetHostToken(gameServer.HostId)
	if hostToken == "" {
		return errors.New("getting host token from db failed")
	}

	channel, ok := events.SSE.GetChannel("/v1/events/" + hostToken)
	if !ok {
		return errors.New("host is not turned on")
	}

	var gameDef game.Game
	err := json.Unmarshal([]byte(gameServer.GameJson), &gameDef)
	if err != nil {
		return errors.New("error when getting game")
	}

	taskMsg := task.Message{
		TaskId:       task.AGENT_START_GS,
		GameServerID: gsId,
		Game:         gameDef,
	}

	taskMsgStr, err := taskMsg.Serialize()
	if err != nil {
		return err
	}

	channel.SendMessage(sse.NewMessage("", taskMsgStr, taskMsg.TaskId.StringValue()))

	model.SetGsStart(model.Redis, gameServer.ID, 1)

	return nil
}

func Stop(gsId string, stopType uint) error {
	gameServer := model.GetGameServer(gsId)
	if gameServer == nil {
		return errors.New("error when getting game server from db")
	}

	hostToken := model.GetHostToken(gameServer.HostId)
	if hostToken == "" {
		return errors.New("error when getting host token from db")
	}

	channel, ok := events.SSE.GetChannel("/v1/events/" + hostToken)
	if !ok {
		return errors.New("host is not turned on")
	}

	taskMsg := task.Message{
		GameServerID: gsId,
	}
	if stopType == 1 {
		taskMsg.TaskId = task.GAME_STOP_SIGKILL
	} else if stopType == 2 {
		taskMsg.TaskId = task.GAME_STOP_SIGTERM
	}

	taskMsgStr, err := taskMsg.Serialize()
	if err != nil {
		return err
	}

	channel.SendMessage(sse.NewMessage("", taskMsgStr, taskMsg.TaskId.StringValue()))

	return nil
}

func Install(gsId string) error {
	gameServer := model.GetGameServer(gsId)
	if gameServer == nil {
		return errors.New("error when getting game server from db")
	}

	hostToken := model.GetHostToken(gameServer.HostId)
	if hostToken == "" {
		return errors.New("error when getting host token from db")
	}

	gameFile := model.GetGameFileByGameIdAndVersion(gameServer.GameId, gameServer.GameVersion)
	if gameFile == nil {
		return errors.New("error when getting game file from db")
	}

	var g game.Game
	err := json.Unmarshal([]byte(gameServer.GameJson), &g)
	if err != nil {
		logrus.Error(err)
	}
	g.DownloadUrl = gameFile.DownloadUrl

	channel, ok := events.SSE.GetChannel("/v1/events/" + hostToken)
	if !ok {
		return errors.New("host is not turned on")
	}

	taskMsg := task.Message{
		TaskId:       task.AGENT_INSTALL_GS,
		Game:         g,
		GameServerID: gsId,
	}

	taskMsgStr, err := taskMsg.Serialize()
	if err != nil {
		return err
	}

	channel.SendMessage(sse.NewMessage("", taskMsgStr, taskMsg.TaskId.StringValue()))

	return nil
}

func SendCommand(gsId string, command string) error {
	gameServer := model.GetGameServer(gsId)
	if gameServer == nil {
		return &lib.Error{ErrorCode: ecode.GsNotFound}
	}

	hostToken := model.GetHostToken(gameServer.HostId)
	if hostToken == "" {
		return errors.New("error when getting host token from db")
	}

	channel, ok := events.SSE.GetChannel("/v1/events/" + hostToken)
	if !ok {
		return errors.New("game server is not turned on")
	}

	taskMsg := task.Message{
		TaskId:       task.GAME_COMMAND,
		GameServerID: gsId,
		Body:         command,
	}

	taskMsgStr, err := taskMsg.Serialize()
	if err != nil {
		return err
	}

	channel.SendMessage(sse.NewMessage("", taskMsgStr, taskMsg.TaskId.StringValue()))

	return nil
}

func Restart(gsId string, stopType uint) error {
	gameServer := model.GetGameServer(gsId)
	if gameServer == nil {
		return &lib.Error{ErrorCode: ecode.GsNotFound}
	}

	hostToken := model.GetHostToken(gameServer.HostId)
	if hostToken == "" {
		return errors.New("error when getting host token from db")
	}

	model.SetGsRestart(model.Redis, gameServer.ID, 0)

	channel, ok := events.SSE.GetChannel("/v1/events/" + hostToken)
	if !ok {
		return errors.New("game server is not turned on")
	}

	taskMsg := task.Message{
		GameServerID: gsId,
	}
	if stopType == 1 {
		taskMsg.TaskId = task.GAME_STOP_SIGKILL
	} else if stopType == 2 {
		taskMsg.TaskId = task.GAME_STOP_SIGTERM
	}

	taskMsgStr, err := taskMsg.Serialize()
	if err != nil {
		return err
	}

	channel.SendMessage(sse.NewMessage("", taskMsgStr, taskMsg.TaskId.StringValue()))

	model.SetGsRestart(model.Redis, gameServer.ID, 1)

	go func() {
		<-time.After(time.Duration(gameServer.StopTimeout) * time.Second)

		val, err := model.GetGsRestart(model.Redis, gameServer.ID)
		if err != nil {
			return
		}

		if val == 1 {
			model.SetGsRestart(model.Redis, gameServer.ID, -1)
		}
	}()

	return nil
}

func Remove(gsId string) error {
	gameServer := model.GetGameServer(gsId)
	if gameServer == nil {
		return &lib.Error{ErrorCode: ecode.GsNotFound}
	}

	hostToken := model.GetHostToken(gameServer.HostId)
	if hostToken == "" {
		return &lib.Error{ErrorCode: ecode.GameNotFound}
	}

	if gameServer.State == 1 {
		taskMsg := task.Message{
			GameServerID: gsId,
			TaskId:       task.GAME_STOP_SIGKILL,
		}
		taskMsgStr, err := taskMsg.Serialize()
		if err != nil {
			return err
		}

		channel, ok := events.SSE.GetChannel("/v1/events/" + hostToken)
		if !ok {
			return errors.New("game server is not turned on")
		}
		channel.SendMessage(sse.NewMessage("", taskMsgStr, taskMsg.TaskId.StringValue()))
	}

	taskMsg := task.Message{
		GameServerID: gsId,
		TaskId:       task.AGENT_REMOVE_GS,
	}
	taskMsgStr, err := taskMsg.Serialize()
	if err != nil {
		return err
	}

	channel, ok := events.SSE.GetChannel("/v1/events/" + hostToken)
	if !ok {
		return errors.New("host is not turned on")
	}
	channel.SendMessage(sse.NewMessage("", taskMsgStr, taskMsg.TaskId.StringValue()))

	permissions := model.GetPermisionsByEndpointRegex("/v1/host/" + gameServer.HostId + "/server/" + gsId + "%")
	for _, perm := range permissions {
		err := model.Delete(perm.ID, perm.Rev)
		if err != nil {
			return err
		}
	}

	err = model.Delete(gameServer.ID, gameServer.Rev)
	if err != nil {
		return err
	}

	return nil
}

func Update(hostId string) error {
	hostToken := model.GetHostToken(hostId)
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

func FileList(gsId string) (*filemanager.Node, error) {
	gameServer := model.GetGameServer(gsId)
	if gameServer == nil {
		return nil, errors.New("error when getting game server from db")
	}

	hostToken := model.GetHostToken(gameServer.HostId)
	if hostToken == "" {
		return nil, errors.New("error when getting host token from db")
	}

	channel, ok := events.SSE.GetChannel("/v1/events/" + hostToken)
	if !ok {
		return nil, errors.New("host is not turned on")
	}

	taskMsg := task.Message{
		TaskId:       task.AGENT_FILE_LIST_GS,
		GameServerID: gsId,
	}

	taskMsgStr, err := taskMsg.Serialize()
	if err != nil {
		return nil, err
	}

	channel.SendMessage(sse.NewMessage("", taskMsgStr, taskMsg.TaskId.StringValue()))

	//wait for files
	pubsub, err := model.GsFilesSubscribe(model.Redis, gsId)
	if err != nil {
		return nil, err
	}
	ch := pubsub.Channel()

	var filesStr string
	for msg := range ch {
		filesStr = msg.Payload
		break
	}
	_ = pubsub.Close()

	var files filemanager.Node
	err = json.Unmarshal([]byte(filesStr), &files)
	if err != nil {
		return nil, err
	}

	return &files, err
}
