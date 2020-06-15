package fetch

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"io"
	"os"
	"path/filepath"

	log "github.com/victron/simpleLogger"
)

// compressing dir in parent location and returning fileName
func Compress(dirName string) (string, error) {
	dir, err := os.Open(dirName)
	if err != nil {
		return "", err
	}
	defer dir.Close()

	dirInfo, err := dir.Stat()
	if err != nil {
		return "", err
	}
	if !dirInfo.IsDir() {
		return "", errors.New("not a dir: " + "dirName")
	}

	files, err := dir.Readdir(0)
	if err != nil {
		return "", err
	}

	if len(files) == 0 {
		log.Warning.Println("dir=", dirName, "is empty")
		return "", errors.New("dir is empty")
	}

	gzFileName := dirName + ".tgz"
	tgzfile, err := os.Create(gzFileName)
	if err != nil {
		return "", err
	}
	defer tgzfile.Close()

	var fileWriter io.WriteCloser = tgzfile
	archiver := gzip.NewWriter(fileWriter)
	archiver.Name = gzFileName
	defer archiver.Close()

	tarfileWriter := tar.NewWriter(archiver)
	defer tarfileWriter.Close()

	for _, fileInfo := range files {
		if fileInfo.IsDir() {
			return "", errors.New("unexpected dir in" + dirName)
		}
		file, err := os.Open(filepath.Join(dir.Name(), fileInfo.Name()))
		if err != nil {
			return "", err
		}
		defer file.Close()

		header, err := tar.FileInfoHeader(fileInfo, fileInfo.Name())
		if err != nil {
			return "", err
		}
		err = tarfileWriter.WriteHeader(header)
		if err != nil {
			return "", err
		}

		_, err = io.Copy(tarfileWriter, file)
		if err != nil {
			return "", err
		}
	}
	return gzFileName, nil
}
