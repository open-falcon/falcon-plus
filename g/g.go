package g

import (
	"runtime"
)

// changelog:
// 1.3.2: add config file `graph_bachends.txt`

// TODO: mv graph cluster config to cfg.json

const (
	VERSION = "1.3.2"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}
