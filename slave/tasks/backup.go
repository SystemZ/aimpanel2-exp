package tasks

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
)

// https://stackoverflow.com/q/13611100/1351857

func TarGzWrite(_path string, tw *tar.Writer, fi os.FileInfo) {
	fr, err := os.Open(_path)
	handleBackupError(err)
	defer fr.Close()

	h := new(tar.Header)
	// we use preserve whole path instead of just filename
	//h.Name = fi.Name()
	h.Name = _path

	h.Size = fi.Size()
	h.Mode = int64(fi.Mode())
	h.ModTime = fi.ModTime()

	err = tw.WriteHeader(h)
	handleBackupError(err)

	_, err = io.Copy(tw, fr)
	handleBackupError(err)
}

func IterDirectory(dirPath string, tw *tar.Writer) {
	dir, err := os.Open(dirPath)
	handleBackupError(err)
	defer dir.Close()
	fis, err := dir.Readdir(0)
	handleBackupError(err)
	for _, fi := range fis {
		curPath := dirPath + "/" + fi.Name()
		if fi.IsDir() {
			//TarGzWrite( curPath, tw, fi )
			IterDirectory(curPath, tw)
		} else {
			fmt.Printf("adding... %s\n", curPath)
			TarGzWrite(curPath, tw, fi)
		}
	}
}

func TarGz(outFilePath string, inPath string) {
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

	IterDirectory(inPath, tw)

	fmt.Println("tar.gz ok")
}

func handleBackupError(err error) {
	if err != nil {
		panic("oh no")
	}
}
