package cron

import (
	"github.com/open-falcon/agent/funcs"
	"time"
)

func InitDataHistory() {
	for {
		funcs.UpdateCpuStat()
		funcs.UpdateDiskStats()
		time.Sleep(time.Second)
	}
}
