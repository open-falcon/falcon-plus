package g

import (
	"log"
	"runtime"
)

// Change Logs
// 0.5.2 use rrdlite other than rrdtool, fix data lost when query failed
//		 rollback to central lib, add filer for debug
const (
	VERSION         = "0.5.2"
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
