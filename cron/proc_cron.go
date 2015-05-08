package cron

import (
	"github.com/open-falcon/transfer/proc"
	"github.com/open-falcon/transfer/sender"
	"github.com/toolkits/container/list"
	"time"
)

const (
	DefaultProcCronPeriod = time.Duration(1) * time.Second //ProcCron的周期,默认1s
)

func StartProcCron() {
	duration := DefaultProcCronPeriod
	for {
		time.Sleep(duration)
		// send cache statistics
		refreshSendingCacheSize()
	}
}

// sending cache
func refreshSendingCacheSize() {
	// total cache size
	proc.JudgeQueuesCnt.Set(calcSendCacheSize(sender.JudgeQueues))
	proc.GraphQueuesCnt.Set(calcSendCacheSize(sender.GraphQueues))
	proc.GraphMigratingQueuesCnt.Set(calcSendCacheSize(sender.GraphMigratingQueues))

	// cache size per node
	for node, list := range sender.JudgeQueues {
		cnt := proc.JudgeQueuesCntPerNode[node]
		if cnt != nil {
			cnt.Set(int64(list.Len()))
		}
	}
	for node, list := range sender.GraphQueues {
		cnt := proc.GraphQueuesCntPerNode[node]
		if cnt != nil {
			cnt.Set(int64(list.Len()))
		}
	}
	for node, list := range sender.GraphMigratingQueues {
		cnt := proc.GraphMigratingQueuesCntPerNode[node]
		if cnt != nil {
			cnt.Set(int64(list.Len()))
		}
	}
}

func calcSendCacheSize(mapList map[string]*list.SafeLinkedListLimited) int64 {
	var cnt int64 = 0
	for _, list := range mapList {
		if list != nil {
			cnt += int64(list.Len())
		}
	}
	return cnt
}
