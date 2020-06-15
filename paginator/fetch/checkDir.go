package fetch

import (
	"errors"
	"os"
	"path/filepath"

	log "github.com/victron/simpleLogger"
)

// test if file type peresent in dir
// return number of such files
func IsFilePresent(fileEnding, dirName string) (int, error) {
	dir, err := os.Open(dirName)
	if err != nil {
		return 0, err
	}
	defer dir.Close()

	dirInfo, err := dir.Stat()
	if err != nil {
		return 0, err
	}
	if !dirInfo.IsDir() {
		return 0, errors.New("not a dir: " + "dirName")
	}

	files, err := dir.Readdir(0)
	if err != nil {
		return 0, err
	}

	if len(files) == 0 {
		log.Warning.Println("dir=", dirName, "is empty")
		return 0, errors.New("dir is empty")
	}

	count := 0
	for _, fileInfo := range files {
		if fileEnding == filepath.Ext(fileInfo.Name()) {
			count++
		}
	}
	return count, nil
}
