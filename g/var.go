package g

import (
	"github.com/toolkits/net"
	"log"
	"sync"
)

var LocalIps []string

func InitVars() {
	var err error
	LocalIps, err = net.IntranetIP()
	if err != nil {
		log.Fatalln("get intranet ip fail:", err)
	}
}

var (
	reportPorts     []int64
	reportPortsLock = new(sync.RWMutex)
)

func ReportPorts() []int64 {
	reportPortsLock.RLock()
	defer reportPortsLock.RUnlock()
	sz := len(reportPorts)
	theClone := make([]int64, sz)
	for i := 0; i < sz; i++ {
		theClone[i] = reportPorts[i]
	}
	return theClone
}

var (
	// tags => {1=>name, 2=>cmdline}
	// e.g. 'name=falcon-agent'=>{1=>falcon-agent}
	// e.g. 'cmdline=xx'=>{2=>xx}
	reportProcs     map[string]map[int]string
	reportProcsLock = new(sync.RWMutex)
)

func ReportProcs() map[string]map[int]string {
	reportProcsLock.RLock()
	defer reportProcsLock.RUnlock()
	sz := len(reportProcs)
	theClone := make(map[string]map[int]string, sz)
	for k, v := range reportProcs {
		theClone[k] = v
	}
	return theClone
}
