package main

import (
	"flag"

	log "github.com/victron/simpleLogger"
)

var user, password *string
var wait_timer *int

func init() {
	user = flag.String("u", "", "user name for authentication")
	password = flag.String("p", "", "user password")
	wait_timer = flag.Int("w", 30, "wait timer befors every request")
	log.FlagsParse()

	// flag.Parse()
}
