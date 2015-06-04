package proc

import (
	nproc "github.com/niean/gotools/proc"
	"log"
	"time"
)

// 索引更新
var (
	IndexUpdateAllCnt = nproc.NewSCounterQps("IndexUpdateAllCnt")
	IndexDeleteCnt    = nproc.NewSCounterQps("IndexDeleteCnt")
)

// 监控数据采集
var (
	CollectorCronCnt = nproc.NewSCounterQps("CollectorCronCnt")
)

// 监控
var (
	MonitorCronCnt            = nproc.NewSCounterQps("MonitorCronCnt")
	MonitorConcurrentErrorCnt = nproc.NewSCounterQps("MonitorConcurrentErrorCnt")
	MonitorAlarmMailCnt       = nproc.NewSCounterQps("MonitorAlarmMailCnt")
)

func Start() {
	log.Println("proc:Start, ok")
}

func GetAll() []interface{} {
	ret := make([]interface{}, 0)

	// index
	ret = append(ret, IndexUpdateAllCnt.Get())
	ret = append(ret, IndexDeleteCnt.Get())

	// collector
	ret = append(ret, CollectorCronCnt.Get())

	// monitor
	ret = append(ret, MonitorCronCnt.Get())
	ret = append(ret, MonitorConcurrentErrorCnt.Get())
	ret = append(ret, MonitorAlarmMailCnt.Get())

	return ret
}

// TODO 临时放在这里了, 考虑放到合适的模块
func FmtUnixTs(ts int64) string {
	return time.Unix(ts, 0).Format("2006-01-02 15:04:05")
}
