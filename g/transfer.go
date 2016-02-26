package g

import (
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/open-falcon/common/model"
)

var (
	TransferLock    *sync.RWMutex                   = new(sync.RWMutex)
	TransferClients map[string]*SingleConnRpcClient = map[string]*SingleConnRpcClient{}
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
		closeTransferClient(addr)
	}
}

func initTransferClient(addr string) {
	TransferLock.Lock()
	defer TransferLock.Unlock()
	TransferClients[addr] = &SingleConnRpcClient{
		RpcServer: addr,
		Timeout:   time.Duration(Config().Transfer.Timeout) * time.Millisecond,
	}
}

func closeTransferClient(addr string) {
	TransferLock.Lock()
	defer TransferLock.Unlock()
	TransferClients[addr].close()
	delete(TransferClients, addr)
}

func updateMetrics(addr string, metrics []*model.MetricValue, resp *model.TransferResponse) bool {
	TransferLock.RLock()
	defer TransferLock.RUnlock()
	err := TransferClients[addr].Call("Transfer.Update", metrics, resp)
	if err != nil {
		log.Println("call Transfer.Update fail", addr, err)
		return false
	}
	return true
}
