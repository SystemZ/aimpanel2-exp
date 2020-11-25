package gameserver

import (
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/lib"
	"gitlab.com/systemz/aimpanel2/lib/ecode"
	"gitlab.com/systemz/aimpanel2/lib/filemanager"
	"gitlab.com/systemz/aimpanel2/lib/task"
	"gitlab.com/systemz/aimpanel2/master/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sync"
	"time"
)

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

func FileRemove(gsId primitive.ObjectID, path string, user model.User) error {
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
		TaskId:       task.AGENT_FILE_REMOVE_GS,
		GameServerID: gsId.Hex(),
		Body:         path,
	}

	err = model.SendTaskToSlave(host.ID, user, taskMsg)
	if err != nil {
		return &lib.Error{ErrorCode: ecode.DbSave}
	}

	return nil
}

func FileServer(gsId primitive.ObjectID, user model.User) (int, error) {
	port, err := lib.GetFreePort()
	if err != nil {
		return 0, err
	}

	gameServer, err := model.GetGameServerById(gsId)
	if err != nil {
		return 0, err
	}

	if gameServer == nil {
		return 0, &lib.Error{ErrorCode: ecode.GsNotFound}
	}

	host, err := model.GetHostById(gameServer.HostId)
	if err != nil {
		return 0, &lib.Error{ErrorCode: ecode.HostNotFound}
	}

	domain, err := model.GetCertDomainByName(host.Domain)
	if err != nil {
		return 0, &lib.Error{ErrorCode: ecode.DomainNotFound}
	}

	cert, err := model.GetCertByDomainId(domain.ID)
	if err != nil {
		return 0, &lib.Error{ErrorCode: ecode.CertNotFound}
	}

	taskMsg := task.Message{
		TaskId:       task.AGENT_FILE_SERVER,
		GameServerID: gsId.Hex(),
		Port:         port,
		Cert:         cert.Cert,
		PrivateKey:   cert.PrivateKey,
	}

	err = model.SendTaskToSlave(host.ID, user, taskMsg)
	if err != nil {
		return 0, &lib.Error{ErrorCode: ecode.DbSave}
	}

	return port, nil
}
