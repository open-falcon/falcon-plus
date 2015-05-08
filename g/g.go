package g

import (
	"runtime"
)

const (
	VERSION = "1.3.1"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}
