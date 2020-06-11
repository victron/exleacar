package paginator

import (
	"net/http"
	"strings"
	"time"

	colly "github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/queue"

	log "github.com/victron/simpleLogger"
)

// hook to star walk on search
func SearchWalker(cookies []*http.Cookie) {
	// mClient := new(mongoClient)
	// if err := (*mClient).Connect(MONGO_LOCAL); err != nil {
	// 	log.Error.Fatal(err)
	// }
	// defer (*mClient).Close()

	q, _ := queue.New(
		1, // Number of consumer threads
		&queue.InMemoryQueueStorage{MaxSize: 10000}, // Use default queue storage
	)

	searchCl := colly.NewCollector(
		colly.AllowedDomains(ALLOWED_DOMAINS...),
		colly.UserAgent(USER_AGENT),
		// colly.CacheDir(CACHE_DIR), // don't forget, this will not use delay :))))
		colly.IgnoreRobotsTxt(),
		// //colly.MaxDepth(2),
		// colly.Async(false), // some problem, not doing requests
	)

	// searchCl := commonCl.Clone()

	searchCl.Limit(&colly.LimitRule{
		// Filter domains affected by this rule
		DomainGlob: "*",
		// Set a delay between requests to these domains
		Delay: time.Duration(*wait_timer) * time.Second,
		// Add an additional random delay
		RandomDelay: time.Duration(*wait_timer) * time.Second,
		Parallelism: 1,
	})

	searchCl.OnRequest(func(r *colly.Request) {
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

	searchCl.OnHTML("a.pagination-next[href]", func(e *colly.HTMLElement) {
		foundURL := e.Request.AbsoluteURL(e.Attr("href"))
		log.Debug.Println("foundURL=", foundURL)
		if err := q.AddURL(foundURL); err != nil {
			log.Error.Fatalln(err)
		}
	})

	searchCl.OnHTML(AUCTION_BLOCK, func(e *colly.HTMLElement) {
		log.Debug.Println("auction block found")
		e.ForEach(`a[href]`, func(_ int, e *colly.HTMLElement) {
			foundURL := e.Request.AbsoluteURL(e.Attr("href"))
			if strings.HasPrefix(foundURL, DETAILS_PREFIX) {
				log.Debug.Println("details foundURL=", foundURL)
				car := new(Car)
				var err error
				(*car).Meta.Url = foundURL
				if (*car).Meta.Id, err = ParseId(foundURL); err != nil {
					log.Error.Fatal(err)
				}
				(*car).Meta.Mdate = time.Now()
				// ProductCollect(productCollector, foundURL)
				// if err := (*car).SaveId(mClient); err != nil {
				// 	log.Error.Fatalln(err)
				// }
				// TODO: add to queue
				// if err := q.AddURL(foundURL); err != nil {
				// 	log.Error.Fatalln(err)
				// }
			}
		})
	})
	searchCl.SetCookies(START_URL, cookies)
	if err := q.AddURL(START_URL); err != nil {
		log.Error.Fatalln(err)
	}

	if err := q.Run(searchCl); err != nil {
		log.Error.Fatalln(err)
	}
}
