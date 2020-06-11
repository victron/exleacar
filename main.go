package main

import (
	"github.com/victron/exleacar/auth"
	"github.com/victron/exleacar/paginator"
	log "github.com/victron/simpleLogger"
)

func main() {
	if *user == "" || *password == "" {
		log.Error.Fatal("user name and password mandatory")
		return
	}
	cookies, err := auth.RecieveCookies("https://www.exleasingcar.com/en",
		*user, *password)
	if err != nil {
		log.Error.Fatalln(err)
	}

	paginator.SearchWalker(cookies)
}
