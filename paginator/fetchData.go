package paginator

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/victron/exleacar/paginator/fetch"
	log "github.com/victron/simpleLogger"
)

// colecting reports and photos
func (car *Car) FetchData(cookies []*http.Cookie, photoNum int) error {
	date := time.Now().Format("2006-01-02")
	dir := filepath.Join(*DATA_DIR, date, strconv.Itoa((*car).Id))
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Error.Fatal(err, "can't create dir=", dir)
	}

	// fetching reports
	if err := car.FetchReports(dir, cookies); err != nil {
		return err
	}

	// check if report present, if not download photos
	if n, _ := fetch.IsFilePresent(".pdf", dir); n == 0 || photoNum > 0 {
		log.Info.Println("getting images")
		// fetching photos
		if err := car.FetchPhotos(dir, cookies, photoNum); err != nil {
			return err
		}
	}

	// compress dir
	archName, err := fetch.Compress(dir)
	if err != nil {
		return err
	} else {
		if err := os.RemoveAll(dir); err != nil {
			return nil
		}

	}
	// update meta data
	(*car).Meta.Fdate = time.Now()
	(*car).Meta.Dir = archName
	(*car).Meta.Fetched = true
	return nil
}

// fetching reports
func (car *Car) FetchReports(dir string, cookies []*http.Cookie) error {
	// var errReturn error
	for n, ldescr := range (*car).Data.Damage {
		if ldescr.Link == "" {
			log.Warning.Println(`no damage reports for id=`, strconv.Itoa((*car).Id), "in link n=", n)
			continue
		}

		// TODO: to remove this part
		// if strings.HasSuffix(ldescr.Link, "/") {
		// 	log.Warning.Println(`found "/" at the end of link=`, ldescr.Link)
		// }

		fileName := filepath.Join(dir, strconv.Itoa(n)+"_report")

		if err := fetch.DownloadFile(fileName, ldescr.Link, cookies); err != nil {
			log.Error.Println(err)
			// errReturn = err
			continue
		}

		if _, err := fetch.RenameFile(fileName); err != nil {
			log.Error.Println("renaming err=", err)
		}

	}
	return nil
}

// fetching photos
// number of photos not less then  MAX_PHOTOS_NUMBER and photoNum
func (car *Car) FetchPhotos(dir string, cookies []*http.Cookie, photoNum int) error {
	for n, photoLink := range (*car).Data.Photos {
		fileName := filepath.Join(dir, strconv.Itoa(n)+"_photo")

		if err := fetch.DownloadFile(fileName, photoLink, cookies); err != nil {
			if strings.HasSuffix(fmt.Sprint(err), "no space left on device") {
				log.Error.Fatalln(err)
			}
			log.Error.Println(err)
			continue
		}
		if _, err := fetch.RenameFile(fileName); err != nil {
			log.Error.Println("renaming err=", err)
		}

		if n > MAX_PHOTOS_NUMBER && n > photoNum {
			break
		}

		time.Sleep(3 * time.Second)
	}
	return nil
}
