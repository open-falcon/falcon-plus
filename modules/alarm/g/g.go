package g

import (
	"log"
	"runtime"
)

const (
	VERSION = "0.2.0"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}
