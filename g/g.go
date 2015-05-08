package g

import (
	"log"
	"runtime"
)

// change log:
// 1.0.7: code refactor for open source
// 1.0.8: bugfix loop init cache
const (
	VERSION = "1.0.8"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}
