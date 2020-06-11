package paginator

import (
	"strings"
	"time"

	colly "github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/queue"

	log "github.com/victron/simpleLogger"
)

// hook to star walk on search
func SearchWalker() {
	log.Debug.Println("in SearchWalker")
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
		colly.CacheDir(CACHE_DIR),
		colly.IgnoreRobotsTxt(),
		// //colly.MaxDepth(2),
		// colly.Async(false), // some problem, not doing requests
	)

	// searchCl := commonCl.Clone()

	searchCl.Limit(&colly.LimitRule{
		// Filter domains affected by this rule
		DomainGlob: "*",
		// Set a delay between requests to these domains
		Delay: NEXT_PAGE * time.Second,
		// Add an additional random delay
		RandomDelay: NEXT_PAGE * time.Second,
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
	// 	log.Debug.Println("URL res.StatusCode=", res.Body)
	// 	// log.Debug.Println("URL from qtx=", res.Request.Ctx.Get("origURL"))
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

	if err := q.AddURL(START_URL); err != nil {
		log.Error.Fatalln(err)
	}

	if err := q.Run(searchCl); err != nil {
		log.Error.Fatalln(err)
	}
}
