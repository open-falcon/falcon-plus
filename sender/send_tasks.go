package sender

import (
	"log"
	"time"

	cmodel "github.com/open-falcon/common/model"
	cutils "github.com/open-falcon/common/utils"
	nsema "github.com/toolkits/concurrent/semaphore"
	nlist "github.com/toolkits/container/list"

	"github.com/open-falcon/gateway/g"
	"github.com/open-falcon/gateway/proc"
)

// send
const (
	DefaultSendTaskSleepInterval = time.Millisecond * 50 //默认睡眠间隔为50ms
)

func startSendTasks() {
	cfg := g.Config()
	concurrent := cfg.Transfer.MaxIdle
	go forward2TransferTask(SenderQueue, concurrent)
}

func forward2TransferTask(Q *nlist.SafeListLimited, concurrent int32) {
	cfg := g.Config()
	batch := int(cfg.Transfer.Batch)
	addr := cfg.Transfer.Addr
	sema := nsema.NewSemaphore(int(concurrent))

	for {
		items := Q.PopBackBy(batch)
		count := len(items)
		if count == 0 {
			time.Sleep(DefaultSendTaskSleepInterval)
			continue
		}

		transItems := make([]*cmodel.MetricValue, count)
		for i := 0; i < count; i++ {
			transItems[i] = convert(items[i].(*cmodel.MetaData))
		}

		sema.Acquire()
		go func(addr string, transItems []*cmodel.MetricValue, count int) {
			defer sema.Release()

			resp := &g.TransferResp{}
			var err error
			sendOk := false
			for i := 0; i < 3; i++ { //最多重试3次
				err = SenderConnPools.Call(addr, "Transfer.Update", transItems, resp)
				if err == nil {
					sendOk = true
					break
				}
				time.Sleep(time.Millisecond * 10)
			}
			// statistics
			if !sendOk {
				log.Printf("send to transfer %s fail: %v connpool:%v", addr, err, SenderConnPools.Proc())
				proc.SendFailCnt.IncrBy(int64(count))
			} else {
				proc.SendCnt.IncrBy(int64(count))
			}
		}(addr, transItems, count)
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
