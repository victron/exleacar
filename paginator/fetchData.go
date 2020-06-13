package paginator

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/victron/exleacar/paginator/fetch"
	log "github.com/victron/simpleLogger"
)

func (car *Car) FetchData(cookies []*http.Cookie) error {
	date := time.Now().Format("2006-01-02")
	dir := filepath.Join(DATA_DIR, date, strconv.Itoa((*car).Id))
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Error.Fatal(err, "can't create dir=", dir)
	}

	// var errReturn error
	for n, ldescr := range (*car).Data.Damage {
		if ldescr.Link == "" {
			log.Warning.Println(`no damage reports for id=`, strconv.Itoa((*car).Id), "in link n=", n)
			continue
		}
		if strings.HasSuffix(ldescr.Link, "/") {
			log.Warning.Println(`found "/" at the end of link=`, ldescr.Link)
		}
		// path := strings.Split(ldescr.Link, "/")
		// fileName := path[len(path)-1]

		fileName := filepath.Join(dir, strconv.Itoa(n))

		if err := fetch.DownloadFile(fileName, ldescr.Link, cookies); err != nil {
			log.Error.Println(err)
			// errReturn = err
			continue
		}

		if _, err := fetch.RenameFile(fileName); err != nil {
			log.Error.Println("renaming err=", err)
		}

	}
	// update meta data
	(*car).Meta.Fdate = time.Now()
	(*car).Meta.Dir = dir
	(*car).Meta.Fetched = true
	time.Sleep(3 * time.Second)
	return nil
}
