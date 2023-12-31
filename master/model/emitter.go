package model

import (
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/lib/filemanager"
	"strings"
)

// https://gist.github.com/miguelmota/835f8dc6935b98b9fbb77766511477ab
// simple to use shovel for tasks sent via SSE
type Emitter map[string]chan string

var (
	GlobalEmitter Emitter
)

func EmitterInit() {
	GlobalEmitter = Emitter{}
}

func GsFilesPublish(gsId string, files *filemanager.Node) error {
	fileListId := "gs-" + gsId + "-filelist"
	logrus.Infof("publishing file list for %v", fileListId)
	GlobalEmitter[fileListId] <- files.String()
	logrus.Info("files published")
	return nil
}

func GsBackupsPublish(gsId string, files []string) error {
	backupListId := "gs-" + gsId + "-backuplist"
	logrus.Infof("publishing backup list for %v", backupListId)
	GlobalEmitter[backupListId] <- strings.Join(files, ",")
	logrus.Info("backup list published")
	return nil
}
