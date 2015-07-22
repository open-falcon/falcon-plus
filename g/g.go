package g

import (
	"log"
	"runtime"
)

// changelog:
// 0.0.1: init project
// 0.0.4: bugfix: set replicas before add node
// 0.0.8: change receiver, mv proc cron to proc pkg, add readme, add gitversion, add config reload, add trace tools
// 0.0.9: fix bugs of conn pool(use transfer's private conn pool, named & minimum)
// 0.0.10: use more efficient proc & sema, rm conn_pool status log
// 0.0.11: fix bug: all graphs' traffic delined when one graph broken down, modify retry interval
// 0.0.14: support sending multi copies to graph node, align ts for judge, add filter

const (
	VERSION      = "0.0.14"
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
