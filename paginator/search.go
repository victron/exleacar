package paginator

import (
	"net/http"
	"strings"
	"time"

	colly "github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/queue"

	"github.com/victron/exleacar/paginator/details"
	log "github.com/victron/simpleLogger"
)

// hook to star walk on search
func SearchWalker(cookies []*http.Cookie, collector *colly.Collector) {
	cl := collector.Clone()
	mClient := new(mongoClient)
	if err := (*mClient).Connect(MONGO_LOCAL); err != nil {
		log.Error.Fatal(err)
	}
	defer (*mClient).Close()

	q, _ := queue.New(
		1, // Number of consumer threads
		&queue.InMemoryQueueStorage{MaxSize: 10000}, // Use default queue storage
	)

	cl.OnRequest(func(r *colly.Request) {
		// r.Headers.Set("Host", HOST)
		// r.Headers.Set("Cookie", COOKIE)
		// r.Headers.Set("Cache-Control", "no-cache")
		log.Info.Println("Visiting URL=:", r.URL)
	})

	// searchCl.OnResponse(func(res *colly.Response) {
	// 	log.Debug.Println("URL res.Request.URL=", fmt.Sprintf("%+v", res.Request.URL))
	// 	log.Debug.Println("URL res.StatusCode=", res.StatusCode)
	// 	// log.Debug.Println("URL res.Headers=", res.Headers)
	// })

	cl.OnHTML("a.pagination-next[href]", func(e *colly.HTMLElement) {
		foundURL := e.Request.AbsoluteURL(e.Attr("href"))
		log.Debug.Println("foundURL=", foundURL)
		if err := q.AddURL(foundURL); err != nil {
			log.Error.Fatalln(err)
		}
	})

	cl.OnHTML(AUCTION_BLOCK, func(e *colly.HTMLElement) {
		log.Debug.Println("auction block found")
		e.ForEach(`a[href]`, func(_ int, e *colly.HTMLElement) {
			foundURL := e.Request.AbsoluteURL(e.Attr("href"))
			if strings.HasPrefix(foundURL, DETAILS_PREFIX) {
				log.Debug.Println("details foundURL=", foundURL)
				car := new(Car)
				var err error
				(*car).Meta.Url = foundURL
				if (*car).ParseUrl() != nil {
					log.Error.Fatal(err)
				}
				(*car).Meta.Mdate = time.Now()

				// TODO: make check if it already in db
				if !(*car).IdPresent(mClient, false) {
					// TODO: move to main packae this call (may be???)
					// get car details
					(*car).Data, err = details.GetDetails((*car).Meta.Url, cookies, collector)
					if err != nil {
						log.Warning.Println("error for id=", (*car).Id, err)
					}
					(*car).Meta.Ddate = time.Now()

					// fentching
					if err := car.FetchData(cookies); err != nil {
						log.Error.Println(err)
					}

					// details recieved
					// reports fetched
					// vin present
					if ((*car).Meta.Ddate != time.Time{}) && (*car).Meta.Fetched && (*car).Data.Vin != "" {
						(*car).Meta.Checked = true
					}

					if err := (*car).InsertFullDoc(mClient); err != nil {
						log.Error.Fatalln(err)
					}
				}

			}
		})
	})
	cl.SetCookies(START_URL, cookies)
	if err := q.AddURL(START_URL); err != nil {
		log.Error.Fatalln(err)
	}

	if err := q.Run(cl); err != nil {
		log.Error.Fatalln(err)
	}
}
