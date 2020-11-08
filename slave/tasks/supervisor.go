package tasks

import (
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/slave/config"
	"os"
	"path/filepath"
)

func GsCleanFiles(gsId string) {
	logrus.Infof("Cleaning files for GS ID %v started", gsId)

	//remove all files in gs dir
	files, err := filepath.Glob(filepath.Join(config.GS_DIR, gsId, "*"))
	if err != nil {
		logrus.Error(err)
	}

	for _, file := range files {
		err = os.RemoveAll(file)
		if err != nil {
			logrus.Error(err)
		}
	}

	logrus.Infof("Cleaning files for GS ID %v finished", gsId)
}
