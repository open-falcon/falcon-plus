package sender

import (
	"log"
	"math/rand"
	"time"

	cmodel "github.com/open-falcon/common/model"
	cutils "github.com/open-falcon/common/utils"
	nsema "github.com/toolkits/concurrent/semaphore"
	nlist "github.com/toolkits/container/list"

	"github.com/open-falcon/gateway/g"
	"github.com/open-falcon/gateway/proc"
)

func startSendTasks() {
	cfg := g.Config()
	concurrent := cfg.Transfer.MaxIdle
	go forward2TransferTask(SenderQueue, concurrent)
}

func forward2TransferTask(Q *nlist.SafeListLimited, concurrent int32) {
	cfg := g.Config()
	batch := int(cfg.Transfer.Batch)
	sema := nsema.NewSemaphore(int(concurrent))
	transNum := len(TransferHostnames)

	for {
		items := Q.PopBackBy(batch)
		count := len(items)
		if count == 0 {
			time.Sleep(time.Millisecond * 50)
			continue
		}

		transItems := make([]*cmodel.MetricValue, count)
		for i := 0; i < count; i++ {
			transItems[i] = convert(items[i].(*cmodel.MetaData))
		}

		sema.Acquire()
		go func(transItems []*cmodel.MetricValue, count int) {
			defer sema.Release()
			var err error

			// 随机遍历transfer列表，直到数据发送成功 或者 遍历完;随机遍历，可以缓解慢transfer
			resp := &g.TransferResp{}
			sendOk := false
			rint := rand.Int()
			for i := 0; i < transNum && !sendOk; i++ {
				idx := (i + rint) % transNum
				host := TransferHostnames[idx]
				addr := TransferMap[host]

				// 对每一个tranfer地址，最多做3次尝试
				TryCnt := 3
				j := 0
				for j = 0; j < TryCnt && !sendOk; j++ { //最多重试3次
					err = SenderConnPools.Call(addr, "Transfer.Update", transItems, resp)
					if err == nil {
						sendOk = true
					}
					time.Sleep(time.Millisecond * 10)
				}
				if j == TryCnt {
					log.Printf("try sending to transfer %s:%s fail: %v connpool:%v", host, addr, err, SenderConnPools.M[addr].Proc())
				}

				// statistics
				if sendOk {
					TransferSendCnt[host].IncrBy(int64(count))
				}
			}

			// statistics
			if !sendOk {
				log.Printf("send to transfer fail, connpool:%v", SenderConnPools.Proc())
				proc.SendFailCnt.IncrBy(int64(count))
			} else {
				proc.SendCnt.IncrBy(int64(count))
			}
		}(transItems, count)
	}
}

func convert(v *cmodel.MetaData) *cmodel.MetricValue {
	return &cmodel.MetricValue{
		Metric:    v.Metric,
		Endpoint:  v.Endpoint,
		Timestamp: v.Timestamp,
		Step:      v.Step,
		Type:      v.CounterType,
		Tags:      cutils.SortedTags(v.Tags),
		Value:     v.Value,
	}
}
