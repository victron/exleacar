package details

import (
	"net/http"

	colly "github.com/gocolly/colly/v2"

	log "github.com/victron/simpleLogger"
)

// hook to star walk on search
func GetDetails(link string, cookies []*http.Cookie, collector *colly.Collector) (Data, error) {
	data := new(Data)

	cl := collector.Clone()

	cl.OnRequest(func(r *colly.Request) {
		// r.Headers.Set("Host", HOST)
		// r.Headers.Set("Cookie", COOKIE)
		// r.Headers.Set("Cache-Control", "no-cache")
		log.Info.Println("Visiting URL=:", r.URL)
	})

	cl.OnResponse(func(res *colly.Response) {
		// log.Debug.Println("URL res.Request.URL=", fmt.Sprintf("%+v", res.Request.URL))
		log.Debug.Println("URL res.StatusCode=", res.StatusCode)
		// log.Debug.Println("URL res.Headers=", res.Headers)
	})

	cl.OnHTML("div.auto-specification table", func(e *colly.HTMLElement) {
		log.Debug.Println("\"div.auto-specification table\" found")
		table := make([][]string, 0)
		e.ForEach("tr", func(_ int, e *colly.HTMLElement) {
			var row []string
			e.ForEach("td", func(_ int, e *colly.HTMLElement) {
				text := e.Text
				row = append(row, text)
			})
			table = append(table, row)
		})
		log.Debug.Println("specification table=", table)
		(*data).Specification = table
	})

	cl.OnHTML("div.auto-supplier-info table", func(e *colly.HTMLElement) {
		log.Debug.Println("\"div.auto-supplier-info table\" found")
		table := make([][]string, 0)
		e.ForEach("tr", func(_ int, e *colly.HTMLElement) {
			var row []string
			e.ForEachWithBreak("td", func(_ int, e *colly.HTMLElement) bool {
				text := e.Text
				// break from row for google translation
				if text == "Original text" {
					return false
				}
				row = append(row, text)
				return true
			})
			table = append(table, row)
		})
		// log.Debug.Println("supplier-info table=", table)
		(*data).SupplierInfo = table
	})

	cl.OnHTML("div.damage-block table", func(e *colly.HTMLElement) {
		log.Debug.Println("\"div.damage-block table\" found")
		table := make([][]LinkDescription, 0, 3)
		tableText := make([][]string, 0, 3)
		e.ForEach("tr", func(_ int, e *colly.HTMLElement) {
			var row []LinkDescription
			e.ForEach("td", func(_ int, e *colly.HTMLElement) {
				linkD := LinkDescription{}

				linkD.Name = e.Text
				var err error
				if linkD.Link, err = FindLink(e.Attr("onclick")); err != nil {
					log.Warning.Println("link not found")
				}
				row = append(row, linkD)
			})
			table = append(table, row)
		})
		log.Debug.Println("damage-block table=", table)
		var err error
		if (*data).Damage, err = FilterDamage(table); err != nil {
			log.Warning.Println(err)
		}
		(*data).RawData.Damage = tableText
	})

	cl.OnHTML("div.photo-block", func(e *colly.HTMLElement) {
		photos := make([]string, 0)
		e.ForEach("div.photo-page", func(_ int, e *colly.HTMLElement) {
			e.ForEach(`a[href]`, func(_ int, e *colly.HTMLElement) {
				link := e.Attr("title")
				if _, err := FindLink(link); err != nil {
					log.Warning.Println("link not found in \"title\" atribute, trying \"onclick\"")
					linkOnclick := e.Attr("onclick")
					if linkOnclick, err = FindLink(linkOnclick); err != nil {
						log.Warning.Println("link not found")
					}
					photos = append(photos, linkOnclick)
				}
				photos = append(photos, link)
			})
		})
		// log.Debug.Println("photos=", photos)
		(*data).Photos = photos
	})

	cl.SetCookies(link, cookies)
	cl.Visit(link)

	var err error
	if (*data).Vin, err = data.GetVin(); err != nil {
		log.Warning.Println(err)
	}

	// log.Debug.Println("data=", *data)
	return *data, nil
}
