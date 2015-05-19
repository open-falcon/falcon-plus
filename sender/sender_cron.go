package sender

import (
	"github.com/open-falcon/transfer/proc"
	"github.com/toolkits/container/list"
	"log"
	"strings"
	"time"
)

const (
	DefaultProcCronPeriod = time.Duration(5) * time.Second    //ProcCron的周期,默认1s
	DefaultLogCronPeriod  = time.Duration(3600) * time.Second //LogCron的周期,默认300s
)

// send_cron程序入口
func startSenderCron() {
	go startProcCron()
	go startLogCron()
}

func startProcCron() {
	for {
		time.Sleep(DefaultProcCronPeriod)
		refreshSendingCacheSize()
	}
}

func startLogCron() {
	for {
		time.Sleep(DefaultLogCronPeriod)
		logConnPoolsProc()
	}
}

func refreshSendingCacheSize() {
	proc.JudgeQueuesCnt.SetCnt(calcSendCacheSize(JudgeQueues))
	proc.GraphQueuesCnt.SetCnt(calcSendCacheSize(GraphQueues))
	proc.GraphMigratingQueuesCnt.SetCnt(calcSendCacheSize(GraphMigratingQueues))
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

func logConnPoolsProc() {
	log.Printf("connPools proc: \n%v", strings.Join(GraphConnPools.Proc(), "\n"))
}
