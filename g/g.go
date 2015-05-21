package g

import (
	"log"
	"runtime"
)

const (
	VERSION         = "0.4.8"
	GAUGE           = "GAUGE"
	DERIVE          = "DERIVE"
	COUNTER         = "COUNTER"
	CACHE_TIME      = 1800000 //ms
	FLUSH_DISK_STEP = 1000    //ms
	DEFAULT_STEP    = 60      //s
	MIN_STEP        = 30      //s
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}
