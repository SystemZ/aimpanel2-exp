package model

import (
	"github.com/sirupsen/logrus"
)

// https://gist.github.com/miguelmota/835f8dc6935b98b9fbb77766511477ab
// simple to use shovel for tasks sent via SSE
type Emitter map[string]chan bool

var (
	GlobalEmitter Emitter
)

func EmitterInit() {
	GlobalEmitter = Emitter{}
}

func GsCleanFilesFinished(gsId string, ok bool) error {
	cleanFilesId := "gs-" + gsId + "-cleanfiles"
	logrus.Infof("publishing cleanfiles for %v", cleanFilesId)
	GlobalEmitter[cleanFilesId] <- ok
	logrus.Info("cleanfiles published")
	return nil
}
