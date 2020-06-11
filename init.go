package main

import (
	"flag"

	log "github.com/victron/simpleLogger"
)

var user, password *string

func init() {
	user = flag.String("u", "", "user name for authentication")
	password = flag.String("p", "", "user password")
	log.FlagsParse()

	// flag.Parse()
}
