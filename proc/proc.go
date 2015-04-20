package proc

import (
	P "github.com/open-falcon/model/proc"
	"log"
	"time"
)

// 索引垃圾清理
var (
	IndexDelete    = P.NewSCounterQps("IndexDelete")
	IndexDeleteCnt = P.NewSCounterBase("IndexDeleteCnt")
)

// transfer监控数据采集
var (
	TransferMonitorCronCnt = P.NewSCounterQps("TransferMonitorCronCnt")
)

func Init() {
	log.Println("proc:Init, ok")
}

func GetAll() []interface{} {
	ret := make([]interface{}, 0)

	// index
	ret = append(ret, IndexDelete)
	ret = append(ret, IndexDeleteCnt)
	// transfer
	ret = append(ret, TransferMonitorCronCnt)

	return ret
}

// TODO 临时放在这里了, 考虑放到合适的模块
func FmtUnixTs(ts int64) string {
	return time.Unix(ts, 0).Format("2006-01-02 15:04:05")
}
