package tasks

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/lib"
	"gitlab.com/systemz/aimpanel2/lib/ahttp"
	"gitlab.com/systemz/aimpanel2/lib/filemanager"
	"gitlab.com/systemz/aimpanel2/lib/task"
	"gitlab.com/systemz/aimpanel2/slave/config"
	"gitlab.com/systemz/aimpanel2/slave/model"
	"mime"
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

type FileDownload struct {
	Path string
}

func GsFileServer(taskMsg task.Message) {
	gsId := taskMsg.GameServerID

	logrus.Infof("starting file server for gs %v", gsId)

	cert, err := lib.ParseCertificate(taskMsg.Cert, taskMsg.PrivateKey)
	if err != nil {
		logrus.Warn("failed to parse cert for gs %v", gsId)
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

	http.HandleFunc("/download", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Disposition")
		if r.Method != "POST" && r.Method != "OPTIONS" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		if r.Method == "OPTIONS" {
			return
		}

		var fd FileDownload
		if err := json.NewDecoder(r.Body).Decode(&fd); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		file, err := os.Open(path.Join(config.GS_DIR, gsId, fd.Path))
		if os.IsNotExist(err) {
			http.Error(w, "file not exist", http.StatusBadRequest)
			return
		}
		defer file.Close()

		fileInfo, err := file.Stat()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		contentDisposition := mime.FormatMediaType("attachment", map[string]string{"filename": fileInfo.Name()})
		w.Header().Set("Content-Disposition", contentDisposition)
		w.Header().Set("Content-Type", "application/octet-stream")
		http.ServeContent(w, r, file.Name(), fileInfo.ModTime(), file)
	})

	go func() {
		if err := server.ListenAndServeTLS("", ""); err != nil {
			logrus.Warnf("file server for %s - %v", gsId, err)
		}
	}()

	time.Sleep(15 * time.Minute)

	err = server.Shutdown(context.Background())
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Infof("file server for %s stopped", gsId)
}
