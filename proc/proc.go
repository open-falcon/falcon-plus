package proc

import (
	P "github.com/open-falcon/model/proc"
	"log"
	"time"
)

// 索引更新
var (
	IndexUpdateAllCnt = P.NewSCounterQps("IndexUpdateAllCnt")
	IndexDeleteCnt    = P.NewSCounterQps("IndexDeleteCnt")
)

// 监控数据采集
var (
	MonitorCronCnt = P.NewSCounterQps("MonitorCronCnt")
)

func Start() {
	log.Println("proc:Start, ok")
}

func GetAll() []interface{} {
	ret := make([]interface{}, 0)

	// index
	ret = append(ret, IndexUpdateAllCnt.Get())
	ret = append(ret, IndexDeleteCnt.Get())

	// monitor
	ret = append(ret, MonitorCronCnt.Get())

	return ret
}

// TODO 临时放在这里了, 考虑放到合适的模块
func FmtUnixTs(ts int64) string {
	return time.Unix(ts, 0).Format("2006-01-02 15:04:05")
}
