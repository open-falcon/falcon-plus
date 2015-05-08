package g

import (
	"log"
	"runtime"
)

// changelog:
// 0.0.1: init project
// 0.0.4: bugfix: set replicas before add node
const (
	VERSION      = "0.0.7"
	GAUGE        = "GAUGE"
	COUNTER      = "COUNTER"
	DERIVE       = "DERIVE"
	DEFAULT_STEP = 60
	MIN_STEP     = 30
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}
