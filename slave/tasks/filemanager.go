package tasks

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/lib"
	"gitlab.com/systemz/aimpanel2/lib/ahttp"
	"gitlab.com/systemz/aimpanel2/lib/filemanager"
	"gitlab.com/systemz/aimpanel2/lib/task"
	"gitlab.com/systemz/aimpanel2/slave/config"
	"gitlab.com/systemz/aimpanel2/slave/model"
	"net/http"
	"os"
	"path"
	"time"
)

func GsFileList(gsId string) {
	logrus.Infof("File list for GS ID %v started", gsId)

	node, err := filemanager.NewTree(config.GS_DIR+"/"+gsId, 100, 64)
	if err != nil {
		logrus.Error(err)
	}

	taskMsg := task.Message{
		TaskId:       task.AGENT_FILE_LIST_GS,
		GameServerID: gsId,
		Files:        node,
	}

	_, err = ahttp.SendTaskData("/v1/events/"+config.HOST_TOKEN, config.HW_ID, taskMsg)
	if err != nil {
		logrus.Error(err)
	}

	logrus.Infof("File list for GS ID %v finished", gsId)
}

func GsFileRemove(gsId string, filePath string) {
	err := os.RemoveAll(path.Join(config.GS_DIR, gsId, filePath))
	if err != nil {
		logrus.Warn(err)
	}
}

func GsFileRemoveTrigger(taskMsg task.Message) {
	supervisorTask := task.Message{
		TaskId:       task.SUPERVISOR_REMOVE_FILE_GS,
		GameServerID: taskMsg.GameServerID,
		Body:         taskMsg.Body,
	}

	model.SendTask(config.REDIS_PUB_SUB_SUPERVISOR_CH, supervisorTask)
}

func GsFileServer(taskMsg task.Message) {
	logrus.Infof("starting file server for gs %v", taskMsg.GameServerID)

	cert, err := lib.ParseCertificate(taskMsg.Cert, taskMsg.PrivateKey)
	if err != nil {
		logrus.Warn("failed to parse cert for gs %v", taskMsg.GameServerID)
		return
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", taskMsg.Port),
		Handler: http.DefaultServeMux,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{*cert},
		},
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello world")
	})

	go func() {
		if err := server.ListenAndServeTLS("", ""); err != nil {
			logrus.Warnf("file server for %s - %v", taskMsg.GameServerID, err)
		}
	}()

	time.Sleep(60 * time.Second)

	err = server.Shutdown(context.Background())
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Infof("file server for %s stopped", taskMsg.GameServerID)
}
