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
