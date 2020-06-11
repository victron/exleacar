package paginator

import (
	"flag"
)

var wait_timer *int

func init() {
	wait_timer = flag.Int("w", 30, "wait timer befors every request")
	// flag.Parse()
}
