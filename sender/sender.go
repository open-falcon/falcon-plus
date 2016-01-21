package sender

import (
	"log"

	pfc "github.com/niean/goperfcounter"
	cmodel "github.com/open-falcon/common/model"
	nlist "github.com/toolkits/container/list"
	nproc "github.com/toolkits/proc"

	"github.com/open-falcon/gateway/g"
	cpool "github.com/open-falcon/gateway/sender/conn_pool"
)

const (
	DefaultSendQueueMaxSize = 1024000 //102.4w
)

var (
	SenderQueue     = nlist.NewSafeListLimited(DefaultSendQueueMaxSize)
	SenderConnPools *cpool.SafeRpcConnPools

	TransferMap         = make(map[string]string, 0)
	TransferHostnames   = make([]string, 0)
	TransferSendCnt     = make(map[string]*nproc.SCounterQps, 0)
	TransferSendFailCnt = make(map[string]*nproc.SCounterQps, 0)
)

func Start() {
	initConnPools()
	startSendTasks()
	startSenderCron()
	log.Println("send.Start, ok")
}

func Push2SendQueue(items []*cmodel.MetaData) {
	for _, item := range items {

		// statistics
		pk := item.PK()
		g.RecvDataTrace.Trace(pk, item)
		g.RecvDataFilter.Filter(pk, item.Value, item)

		isOk := SenderQueue.PushFront(item)

		// statistics
		if !isOk {
			pfc.Meter("SendDrop", 1)
		}
	}
}

func initConnPools() {
	cfg := g.Config()

	// init transfer global configs
	addrs := make([]string, 0)
	for hn, addr := range cfg.Transfer.Cluster {
		TransferHostnames = append(TransferHostnames, hn)
		addrs = append(addrs, addr)
		TransferMap[hn] = addr
	}

	// init transfer send cnt
	for hn, addr := range cfg.Transfer.Cluster {
		TransferSendCnt[hn] = nproc.NewSCounterQps(hn + ":" + addr)
		TransferSendFailCnt[hn] = nproc.NewSCounterQps(hn + ":" + addr)
	}

	// init conn pools
	SenderConnPools = cpool.CreateSafeRpcConnPools(cfg.Transfer.MaxConns, cfg.Transfer.MaxIdle,
		cfg.Transfer.ConnTimeout, cfg.Transfer.CallTimeout, addrs)
}
