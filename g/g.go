package g

import (
	"log"
	"runtime"
)

// changelog:
// 0.0.1: init project
// 0.0.3: change conn pools, enlarge sending queue
// 0.0.4: use relative paths in 'import', mv conn_pool to central libs
// 0.0.5: use absolute path in 'import'

const (
	VERSION      = "0.0.5"
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
