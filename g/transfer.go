package g

import (
	"log"
	"sync"
	"time"

	"github.com/CMGS/consistent"
	"github.com/open-falcon/common/model"
)

var (
	Transfers       *consistent.Consistent
	TransferLock    *sync.RWMutex                   = new(sync.RWMutex)
	TransferClients map[string]*SingleConnRpcClient = map[string]*SingleConnRpcClient{}
)

func InitTransfers() {
	Transfers = consistent.New()
	for _, transfer := range Config().Transfer.Addrs {
		Transfers.Add(transfer)
	}
}

func SendMetrics(metrics []*model.MetricValue, resp model.TransferResponse) {
	for offset := 0; offset < Transfers.Len(); offset++ {
		addr, _ := Transfers.Get(metrics[0].Endpoint, offset)
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
	delete(TransferClients, addr)
}

func updateMetrics(addr string, metrics []*model.MetricValue, resp model.TransferResponse) bool {
	TransferLock.RLock()
	defer TransferLock.RUnlock()
	err := TransferClients[addr].Call("Transfer.Update", metrics, &resp)
	if err != nil {
		log.Println("call Transfer.Update fail", addr, err)
		return false
	}
	return true
}
