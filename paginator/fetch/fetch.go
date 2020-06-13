package fetch

import (
	"errors"
	"io"
	"mime"
	"net/http"
	"os"
	"time"

	log "github.com/victron/simpleLogger"
)

func DownloadFile(filepath string, url string, cookies []*http.Cookie) error {

	log.Debug.Println("fetching=", url)
	log.Debug.Println("saving=", filepath)
	client := &http.Client{
		Timeout: time.Minute,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	// TODO: jar??????
	for _, c := range cookies {
		req.AddCookie(c)
	}

	res, err := client.Do(req)
	if err != nil {
		log.Error.Println(err)
	}
	defer res.Body.Close()
	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		log.Warning.Println("status code=", res.StatusCode)
		return errors.New("bad responce")
	}

	// resp, err := client.Get(url)
	// resp, err := http.Get(url)
	// defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, res.Body)
	return err
}

// DetectContent and rename file
func RenameFile(fileName string) (string, error) {
	f, err := os.Open(fileName)
	if err != nil {
		log.Error.Println(err)
		return fileName, err
	}

	buffer := make([]byte, 512)
	_, err = f.Read(buffer)
	if err != nil {
		return fileName, err
	}
	f.Close()

	contentType := http.DetectContentType(buffer)
	fileEndings, err := mime.ExtensionsByType(contentType)
	if err != nil {
		log.Error.Println("file detection problem=", err)
		return fileName, err
	}
	log.Debug.Println("fileType=", fileEndings[0])
	newName := fileName + fileEndings[0]

	if err := os.Rename(fileName, newName); err != nil {
		log.Error.Println(err)
		return fileName, err
	}
	return newName, nil
}
