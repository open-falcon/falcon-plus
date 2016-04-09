package g

import (
	"log"
	"runtime"
)

// changelog:
// 0.0.1: init project
// 0.0.3: add readme, add gitversion, modify proc, add config reload
// 0.0.4: make collector configurable, add monitor cron, adjust index db
// Changes: send turning-ok only after alarm happens, add conn timeout for http
//			maybe fix bug of 'too many open files', rollback to central lib
// 0.0.5: move self.monitor to anteye
// 0.0.6: make index update configurable, use global time formater
// 0.0.7: fix bug of index_update_all
// 0.0.8: add agents' house_keeper, use relative paths in 'import'
// 0.0.9: gen falcon.task.alive, use common module, use absolute paths in import
// 0.0.10: rm monitor, add controller for index cleaner

const (
	VERSION = "0.0.10"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}
