package main

import (
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/victron/exleacar/auth"
	"github.com/victron/exleacar/paginator"
	log "github.com/victron/simpleLogger"
)

func main() {
	if *user == "" || *password == "" || *paginator.DATA_DIR == "" {
		log.Error.Fatal("user name, password, dir mandatory")
		return
	}
	cookies, err := auth.RecieveCookies("https://www.exleasingcar.com/en",
		*user, *password)
	if err != nil {
		log.Error.Fatalln(err)
	}

	Cl := colly.NewCollector(
		colly.AllowedDomains(ALLOWED_DOMAINS...),
		colly.UserAgent(USER_AGENT),
		// colly.CacheDir(CACHE_DIR), // don't forget, this will not use delay :))))
		colly.IgnoreRobotsTxt(),
		// //colly.MaxDepth(2),
		// colly.Async(false), // some problem, not doing requests
	)

	Cl.Limit(&colly.LimitRule{
		// Filter domains affected by this rule
		DomainGlob: "*",
		// Set a delay between requests to these domains
		Delay: time.Duration(*wait_timer) * time.Second,
		// Add an additional random delay
		RandomDelay: time.Duration(*wait_timer) * time.Second,
		Parallelism: 1,
	})

	paginator.SearchWalker(cookies, Cl)

	// test for details
	// details.GetDetails("https://www.exleasingcar.com/en/auto-details/6896925",
	// 	cookies, Cl)
	// details.GetDetails("https://www.exleasingcar.com/en/auto-details/6898242",
	// 	cookies, Cl)
}
