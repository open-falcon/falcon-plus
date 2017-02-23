package config

import (
	"log"
	"runtime"
)

// change log:
const (
	VERSION = "0.0.1"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}
