package gameserver

import (
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/lib"
	"gitlab.com/systemz/aimpanel2/lib/ecode"
	"gitlab.com/systemz/aimpanel2/lib/filemanager"
	"gitlab.com/systemz/aimpanel2/lib/game"
	"gitlab.com/systemz/aimpanel2/lib/task"
	"gitlab.com/systemz/aimpanel2/master/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sync"
	"time"
)

func Start(gsId primitive.ObjectID, user model.User) error {
	gameServer, err := model.GetGameServerById(gsId)
	if err != nil {
		return err
	}
	if gameServer == nil {
		return errors.New("getting game server from db failed")
	}

	host, err := model.GetHostById(gameServer.HostId)
	if err != nil {
		return &lib.Error{ErrorCode: ecode.HostNotFound}
	}

	var gameDef game.Game
	err = json.Unmarshal([]byte(gameServer.GameJson), &gameDef)
	if err != nil {
		return errors.New("error when getting game")
	}

	// FIXME validate if current plan allows custom cmd, throw HTTP error if plan is too low
	gameDef.CustomCommandStart = gameServer.CustomCmdStart

	taskMsg := task.Message{
		TaskId:       task.AGENT_START_GS,
		GameServerID: gsId.Hex(),
		Game:         &gameDef,
	}

	err = model.SendTaskToSlave(host.ID, user, taskMsg)
	if err != nil {
		return &lib.Error{ErrorCode: ecode.DbSave}
	}

	return nil
}

func Stop(gsId primitive.ObjectID, stopType uint, user model.User) error {
	gameServer, err := model.GetGameServerById(gsId)
	if err != nil {
		return err
	}

	if gameServer == nil {
		return errors.New("error when getting game server from db")
	}

	host, err := model.GetHostById(gameServer.HostId)
	if err != nil {
		return &lib.Error{ErrorCode: ecode.HostNotFound}
	}

	taskMsg := task.Message{
		GameServerID: gsId.Hex(),
	}
	if stopType == 1 {
		taskMsg.TaskId = task.GAME_STOP_SIGKILL
	} else if stopType == 2 {
		taskMsg.TaskId = task.GAME_STOP_SIGTERM
	}

	err = model.SendTaskToSlave(host.ID, user, taskMsg)
	if err != nil {
		return &lib.Error{ErrorCode: ecode.DbSave}
	}

	return nil
}

func Install(gsId primitive.ObjectID, user model.User) error {
	gameServer, err := model.GetGameServerById(gsId)
	if err != nil {
		return err
	}

	if gameServer == nil {
		return errors.New("error when getting game server from db")
	}

	host, err := model.GetHostById(gameServer.HostId)
	if err != nil {
		return &lib.Error{ErrorCode: ecode.HostNotFound}
	}

	gameFile, err := model.GetGameFileByGameIdAndVersion(gameServer.GameId, gameServer.GameVersion)
	if err != nil {
		return err
	}

	if gameFile == nil {
		return errors.New("error when getting game file from db")
	}

	var g game.Game
	err = json.Unmarshal([]byte(gameServer.GameJson), &g)
	if err != nil {
		logrus.Error(err)
	}
	g.DownloadUrl = gameFile.DownloadUrl

	taskMsg := task.Message{
		TaskId:       task.AGENT_INSTALL_GS,
		Game:         &g,
		GameServerID: gsId.Hex(),
	}

	err = model.SendTaskToSlave(host.ID, user, taskMsg)
	if err != nil {
		return &lib.Error{ErrorCode: ecode.DbSave}
	}

	return nil
}

func SendCommand(gsId primitive.ObjectID, command string, user model.User) error {
	gameServer, err := model.GetGameServerById(gsId)
	if err != nil {
		return err
	}

	if gameServer == nil {
		return &lib.Error{ErrorCode: ecode.GsNotFound}
	}

	host, err := model.GetHostById(gameServer.HostId)
	if err != nil {
		return &lib.Error{ErrorCode: ecode.HostNotFound}
	}

	taskMsg := task.Message{
		TaskId:       task.GAME_COMMAND,
		GameServerID: gsId.Hex(),
		Body:         command,
	}

	err = model.SendTaskToSlave(host.ID, user, taskMsg)
	if err != nil {
		return &lib.Error{ErrorCode: ecode.DbSave}
	}

	return nil
}

func Restart(gsId primitive.ObjectID, stopType uint, user model.User) error {
	gameServer, err := model.GetGameServerById(gsId)
	if err != nil {
		return err
	}

	if gameServer == nil {
		return &lib.Error{ErrorCode: ecode.GsNotFound}
	}

	host, err := model.GetHostById(gameServer.HostId)
	if err != nil {
		return &lib.Error{ErrorCode: ecode.HostNotFound}
	}

	var gameDef game.Game
	err = json.Unmarshal([]byte(gameServer.GameJson), &gameDef)
	if err != nil {
		return errors.New("error when getting game")
	}

	taskMsg := task.Message{
		TaskId:       task.GAME_RESTART,
		GameServerID: gsId.Hex(),
		StopTimeout:  gameServer.StopTimeout,
		Game:         &gameDef,
	}

	err = model.SendTaskToSlave(host.ID, user, taskMsg)
	if err != nil {
		return &lib.Error{ErrorCode: ecode.DbSave}
	}

	return nil
}

func Remove(gsId primitive.ObjectID, user model.User) error {
	gameServer, err := model.GetGameServerById(gsId)
	if err != nil {
		return err
	}
	if gameServer == nil {
		return &lib.Error{ErrorCode: ecode.GsNotFound}
	}

	host, err := model.GetHostById(gameServer.HostId)
	if err != nil {
		return &lib.Error{ErrorCode: ecode.HostNotFound}
	}

	if gameServer.State == 1 {
		taskMsg := task.Message{
			GameServerID: gsId.Hex(),
			TaskId:       task.GAME_STOP_SIGKILL,
		}

		err = model.SendTaskToSlave(host.ID, user, taskMsg)
		if err != nil {
			return &lib.Error{ErrorCode: ecode.DbSave}
		}
	}

	taskMsg := task.Message{
		GameServerID: gsId.Hex(),
		TaskId:       task.AGENT_REMOVE_GS,
	}

	err = model.SendTaskToSlave(host.ID, user, taskMsg)
	if err != nil {
		return &lib.Error{ErrorCode: ecode.DbSave}
	}

	//permissions, err := model.GetPermisionsByEndpointRegex("/v1/host/" + gameServer.HostId.Hex() + "/server/" + gsId.Hex() + "%")
	//if err != nil {
	//	return err
	//}
	//
	//for _, perm := range permissions {
	//	err := model.Delete(&perm)
	//	if err != nil {
	//		return err
	//	}
	//}

	err = model.Delete(gameServer)
	if err != nil {
		return err
	}

	return nil
}

func FileList(gsId primitive.ObjectID, user model.User) (res *filemanager.Node, err error) {
	gameServer, err := model.GetGameServerById(gsId)
	if err != nil {
		return
	}

	if gameServer == nil {
		return nil, errors.New("error when getting game server from db")
	}

	host, err := model.GetHostById(gameServer.HostId)
	if err != nil {
		return nil, &lib.Error{ErrorCode: ecode.HostNotFound}
	}

	taskMsg := task.Message{
		TaskId:       task.AGENT_FILE_LIST_GS,
		GameServerID: gsId.Hex(),
	}

	// FIXME make safe for concurrent requests
	fileListId := "gs-" + gsId.Hex() + "-filelist"
	logrus.Infof("waiting for %v", fileListId)
	model.GlobalEmitter[fileListId] = make(chan string)

	var wg sync.WaitGroup
	wg.Add(1)
	// cancel after 10 seconds of waiting
	go func() {
		time.Sleep(time.Second * 10)
		return
	}()

	// wait for slave to provide result
	go func() {
		for {
			select {
			case msg := <-model.GlobalEmitter[fileListId]:
				var files filemanager.Node
				err = json.Unmarshal([]byte(msg), &files)
				if err != nil {
					wg.Done()
					return
				}
				res = &files
				wg.Done()
			}
		}
	}()

	// send task to slave
	err = model.SendTaskToSlave(host.ID, user, taskMsg)
	if err != nil {
		err = &lib.Error{ErrorCode: ecode.DbSave}
		wg.Done()
	}

	// wait for result or 10 seconds, whatever comes first
	wg.Wait()
	return
}

func Shutdown(gsId primitive.ObjectID, user model.User) error {
	gameServer, err := model.GetGameServerById(gsId)
	if err != nil {
		return err
	}

	if gameServer == nil {
		return errors.New("error when getting game server from db")
	}

	host, err := model.GetHostById(gameServer.HostId)
	if err != nil {
		return &lib.Error{ErrorCode: ecode.HostNotFound}
	}

	var g game.Game
	err = json.Unmarshal([]byte(gameServer.GameJson), &g)
	if err != nil {
		logrus.Error(err)
	}

	taskMsg := task.Message{
		TaskId:       task.GAME_SHUTDOWN,
		Game:         &g,
		GameServerID: gsId.Hex(),
	}

	err = model.SendTaskToSlave(host.ID, user, taskMsg)
	if err != nil {
		return &lib.Error{ErrorCode: ecode.DbSave}
	}

	return nil
}
