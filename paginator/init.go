package paginator

import "flag"

var DATA_DIR *string

func init() {
	DATA_DIR = flag.String("dir", "", "path to store reports and photos")
}
