package g

import (
	"log"
	"runtime"
)

const (
	VERSION = "2.0.2"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}
