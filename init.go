package main

import (
	log "github.com/victron/simpleLogger"
)

// var wait_timer *int

func init() {
	// wait_timer = flag.Int("w", 30, "wait timer befors every request")
	log.FlagsParse()

	// flag.Parse()
}
