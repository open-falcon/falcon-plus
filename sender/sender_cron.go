package sender

import (
	"time"

	"github.com/open-falcon/gateway/proc"
)

const (
	DefaultProcCronPeriod = time.Duration(5) * time.Second //ProcCron的周期,默认1s
)

// send_cron程序入口
func startSenderCron() {
	go startProcCron()
}

func startProcCron() {
	for {
		time.Sleep(DefaultProcCronPeriod)
		refreshSendingCacheSize()
	}
}

func refreshSendingCacheSize() {
	proc.SendQueuesCnt.SetCnt(int64(SenderQueue.Len()))
}
