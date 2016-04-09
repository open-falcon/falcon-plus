package proc

import (
	nproc "github.com/toolkits/proc"
	"log"
)

// 索引更新
var (
	IndexUpdateCnt      = nproc.NewSCounterQps("IndexUpdateCnt")
	IndexUpdateErrorCnt = nproc.NewSCounterQps("IndexUpdateErrorCnt")
	IndexDeleteCnt      = nproc.NewSCounterQps("IndexDeleteCnt")
)

// 监控数据采集
var (
	CollectorCronCnt = nproc.NewSCounterQps("CollectorCronCnt")
)

func Start() {
	log.Println("proc.Start ok")
}

func GetAll() []interface{} {
	ret := make([]interface{}, 0)

	// index
	ret = append(ret, IndexUpdateCnt.Get())
	ret = append(ret, IndexUpdateErrorCnt.Get())
	ret = append(ret, IndexDeleteCnt.Get())

	// collector
	ret = append(ret, CollectorCronCnt.Get())

	return ret
}
