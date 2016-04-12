package sender

import (
	"time"

	pfc "github.com/niean/goperfcounter"
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
	pfc.Gauge("SendQueueSize", int64(SenderQueue.Len()))
}
