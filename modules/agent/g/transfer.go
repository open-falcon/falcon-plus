package g

import (
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/open-falcon/falcon-plus/common/model"
)

var (
	TransferClientsLock *sync.RWMutex                   = new(sync.RWMutex)
	TransferClients     map[string]*SingleConnRpcClient = map[string]*SingleConnRpcClient{}
)

func SendMetrics(metrics []*model.MetricValue, resp *model.TransferResponse) {
	rand.Seed(time.Now().UnixNano())
	for _, i := range rand.Perm(len(Config().Transfer.Addrs)) {
		addr := Config().Transfer.Addrs[i]
		if _, ok := TransferClients[addr]; !ok {
			initTransferClient(addr)
		}
		if updateMetrics(addr, metrics, resp) {
			break
		}
	}
}

func initTransferClient(addr string) {
	TransferClientsLock.Lock()
	defer TransferClientsLock.Unlock()
	TransferClients[addr] = &SingleConnRpcClient{
		RpcServer: addr,
		Timeout:   time.Duration(Config().Transfer.Timeout) * time.Millisecond,
	}
}

func updateMetrics(addr string, metrics []*model.MetricValue, resp *model.TransferResponse) bool {
	TransferClientsLock.RLock()
	defer TransferClientsLock.RUnlock()
	err := TransferClients[addr].Call("Transfer.Update", metrics, resp)
	if err != nil {
		log.Println("call Transfer.Update fail", addr, err)
		return false
	}
	return true
}
