package tasks

import (
	"archive/tar"
	gzip "github.com/klauspost/pgzip"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"strings"
)

// https://stackoverflow.com/q/13611100/1351857

func TarGzWrite(_path string, tw *tar.Writer, fi os.FileInfo, filter bool, filterPath string) {
	fr, err := os.Open(_path)
	handleBackupError(err)
	defer fr.Close()

	h := new(tar.Header)

	// this is ok for flat tar
	//h.Name = fi.Name()

	if filter {
		// normally whole absolute path is in tar like provided what to backup
		// we should use filter this out for more flexibility while restoring
		h.Name = "files" + strings.TrimPrefix(_path, filterPath)
	} else {
		// we use preserve whole path instead of just filename
		h.Name = _path
	}

	h.Size = fi.Size()
	h.Mode = int64(fi.Mode())
	h.ModTime = fi.ModTime()

	err = tw.WriteHeader(h)
	handleBackupError(err)

	_, err = io.Copy(tw, fr)
	handleBackupError(err)
}

func IterDirectory(dirPath string, tw *tar.Writer, filter bool, filterPath string) {
	dir, err := os.Open(dirPath)
	handleBackupError(err)
	defer dir.Close()
	fis, err := dir.Readdir(0)
	handleBackupError(err)
	for _, fi := range fis {
		curPath := dirPath + "/" + fi.Name()
		if fi.IsDir() {
			//TarGzWrite( curPath, tw, fi )
			IterDirectory(curPath, tw, filter, filterPath)
		} else {
			logrus.Debugf("Adding to tar: %s", curPath)
			TarGzWrite(curPath, tw, fi, filter, filterPath)
		}
	}
}

func TarGz(outFilePath string, inPath string, filterAbsolutePath bool) {
	// file write
	fw, err := os.Create(outFilePath)
	handleBackupError(err)
	defer fw.Close()

	// gzip write
	gw := gzip.NewWriter(fw)
	defer gw.Close()

	// tar write
	tw := tar.NewWriter(gw)
	defer tw.Close()

	IterDirectory(inPath, tw, filterAbsolutePath, inPath)
	logrus.Info("Completed .tar.gz creation")
}

func handleBackupError(err error) {
	if err != nil {
		// FIXME handle errors for backups earlier, remove this func
		panic("Something went wrong with backup")
	}
}
