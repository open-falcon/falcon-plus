package g

import (
	"github.com/toolkits/net"
	"log"
	"sync"
	"time"
)

var LocalIps []string

func InitLocalIps() {
	var err error
	LocalIps, err = net.IntranetIP()
	if err != nil {
		log.Fatalln("get intranet ip fail:", err)
	}
}

var (
	HbsClient      *SingleConnRpcClient
	TransferClient *SingleConnRpcClient
)

func InitRpcClients() {
	if Config().Heartbeat.Enabled {
		HbsClient = &SingleConnRpcClient{
			RpcServer: Config().Heartbeat.Addr,
			Timeout:   time.Duration(Config().Heartbeat.Timeout) * time.Millisecond,
		}
	}

	if Config().Transfer.Enabled {
		TransferClient = &SingleConnRpcClient{
			RpcServer: Config().Transfer.Addr,
			Timeout:   time.Duration(Config().Transfer.Timeout) * time.Millisecond,
		}
	}
}

func SendToTransfer(metrics []*MetricValue) {
	if len(metrics) == 0 {
		return
	}

	debug := Config().Debug

	if debug {
		log.Println("=>", metrics[0])
	}

	var resp TransferResp
	TransferClient.Call("Transfer.Update", metrics, &resp)

	if debug {
		log.Println("<=", &resp)
	}
}

var (
	reportPorts     []int64
	reportPortsLock = new(sync.RWMutex)
)

func ReportPorts() []int64 {
	reportPortsLock.RLock()
	defer reportPortsLock.RUnlock()
	return reportPorts
}

func SetReportPorts(ports []int64) {
	reportPortsLock.Lock()
	defer reportPortsLock.Unlock()
	reportPorts = ports
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
	return reportProcs
}

func SetReportProcs(procs map[string]map[int]string) {
	reportProcsLock.Lock()
	defer reportProcsLock.Unlock()
	reportProcs = procs
}

var (
	whiteIPs    []*WhiteIP
	whiteIpLock = new(sync.Mutex)
)

func WhiteIps() []*WhiteIP {
	whiteIpLock.Lock()
	defer whiteIpLock.Unlock()
	return whiteIPs
}

func SetWhiteIps(ips []*WhiteIP) {
	whiteIpLock.Lock()
	defer whiteIpLock.Unlock()
	whiteIPs = ips
}
