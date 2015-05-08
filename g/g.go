package g

import (
	"log"
	"runtime"
)

// change log
// 2.0.1: bugfix HistoryData limit
const (
	VERSION = "2.0.1"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}
