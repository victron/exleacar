package paginator

import (
	"fmt"
	"net/http"
	"strconv"
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
		// r.Ctx.Put("OnRequest", fmt.Sprintf("%v", r.URL))
		log.Info.Println("Visiting URL=:", r.URL)
	})

	// cl.OnResponse(func(res *colly.Response) {
	// 	log.Debug.Println("URL res.Request.URL=", fmt.Sprintf("%+v", res.Request.URL))
	// 	log.Debug.Println("URL res.StatusCode=", res.StatusCode)
	// 	res.Ctx.Put("OnResponse", fmt.Sprintf("%v", res.Request.URL))
	// 	// log.Debug.Println("URL res.Headers=", res.Headers)
	// })

	// cl.OnHTML("a.pagination-next[href]", func(e *colly.HTMLElement) {
	cl.OnHTML("div.pagination-block", func(e *colly.HTMLElement) {
		log.Debug.Println("found=", "div.pagination-block")
		foundURL := e.ChildAttr(`a.pagination-next`, "href")

		// normal walk on pages
		if foundURL != "" {
			foundURL = e.Request.AbsoluteURL(foundURL)
			log.Debug.Println("foundURL=", foundURL)
			if err := q.AddURL(foundURL); err != nil {
				log.Error.Fatalln(err)
			}
			return
		}

		// addiing to queue search on all
		if foundURL == "" && strings.HasPrefix(e.Request.URL.String(), CUSTOM_SEARCH_URL) {
			log.Info.Println("adding ALL_SEARCH_URL to queue")
			cl.SetCookies(ALL_SEARCH_URL, cookies)
			if err := q.AddURL(ALL_SEARCH_URL); err != nil {
				log.Error.Fatalln(err)
			}
		}
	})

	cl.OnHTML(AUCTION_BLOCK, func(e *colly.HTMLElement) {
		log.Debug.Println("auction block found")
		count := 0
		e.ForEach(`a[href]`, func(_ int, e *colly.HTMLElement) {
			foundURL := e.Request.AbsoluteURL(e.Attr("href"))
			if strings.HasPrefix(foundURL, DETAILS_PREFIX) {
				log.Debug.Println("details foundURL=", foundURL)
				count++
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
					e.Response.Ctx.Put("searchUrl", fmt.Sprintf("%v", e.Request.URL))
					e.Response.Ctx.Put("count", strconv.Itoa(count))
					(*car).Data, err = details.GetDetails((*car).Meta.Url, e.Response.Ctx, cookies, collector)
					if err != nil {
						log.Warning.Println("error for id=", (*car).Id, err)
					}
					(*car).Meta.Ddate = time.Now()

					// fentching
					if err := car.FetchData(cookies); err != nil {
						log.Error.Println(err)
					}

					// TODO: change logic later; when details should be update second time
					// details recieved
					// reports fetched
					// vin present
					// if ((*car).Meta.Ddate != time.Time{}) && (*car).Meta.Fetched && (*car).Data.Vin != "" {
					// 	(*car).Meta.Checked = true
					// }
					(*car).Meta.Checked = true

					if err := (*car).InsertFullDoc(mClient); err != nil {
						log.Error.Fatalln(err)
					}
				}

			}
		})
	})
	cl.SetCookies(CUSTOM_SEARCH_URL, cookies)
	if err := q.AddURL(CUSTOM_SEARCH_URL); err != nil {
		log.Error.Fatalln(err)
	}

	if err := q.Run(cl); err != nil {
		log.Error.Fatalln(err)
	}
}
