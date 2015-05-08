package proc

import (
	P "github.com/open-falcon/common/proc"
	"time"
)

// 索引增量更新
var (
	IndexUpdateIncr         = P.NewSCounterQps("IndexUpdateIncr")
	IndexUpdateIncrCnt      = P.NewSCounterBase("IndexUpdateIncrCnt")
	IndexUpdateIncrErrorCnt = P.NewSCounterQps("IndexUpdateIncrErrorCnt")

	IndexUpdateIncrDbEndpointSelectCnt = P.NewSCounterQps("IndexUpdateIncrDbEndpointSelectCnt")
	IndexUpdateIncrDbEndpointInsertCnt = P.NewSCounterQps("IndexUpdateIncrDbEndpointInsertCnt")

	IndexUpdateIncrDbTagEndpointSelectCnt = P.NewSCounterQps("IndexUpdateIncrDbTagEndpointSelectCnt")
	IndexUpdateIncrDbTagEndpointInsertCnt = P.NewSCounterQps("IndexUpdateIncrDbTagEndpointInsertCnt")

	IndexUpdateIncrDbEndpointCounterSelectCnt = P.NewSCounterQps("IndexUpdateIncrDbEndpointCounterSelectCnt")
	IndexUpdateIncrDbEndpointCounterInsertCnt = P.NewSCounterQps("IndexUpdateIncrDbEndpointCounterInsertCnt")
	IndexUpdateIncrDbEndpointCounterUpdateCnt = P.NewSCounterQps("IndexUpdateIncrDbEndpointCounterUpdateCnt")
)

// 索引全量更新
var (
	IndexUpdateAll         = P.NewSCounterQps("IndexUpdateAll")
	IndexUpdateAllCnt      = P.NewSCounterBase("IndexUpdateAllCnt")
	IndexUpdateAllErrorCnt = P.NewSCounterQps("IndexUpdateAllErrorCnt")
)

// 索引缓存大小
var (
	IndexedItemCacheCnt            = P.NewSCounterBase("IndexedItemCacheCnt")
	UnIndexedItemCacheCnt          = P.NewSCounterBase("UnIndexedItemCacheCnt")
	IndexDbEndpointCacheCnt        = P.NewSCounterBase("IndexDbEndpointCacheCnt")
	IndexDbTagEndpointCacheCnt     = P.NewSCounterBase("IndexDbTagEndpointCacheCnt")
	IndexDbEndpointCounterCacheCnt = P.NewSCounterBase("IndexDbEndpointCounterCacheCnt")
)

// Rpc
var (
	GraphRpcRecvCnt = P.NewSCounterQps("GraphRpcRecvCnt")
)

func GetAll() []interface{} {
	ret := make([]interface{}, 0)

	// rpc recv
	ret = append(ret, GraphRpcRecvCnt.Get())

	// index update all
	ret = append(ret, IndexUpdateAll.Get())
	ret = append(ret, IndexUpdateAllCnt.Get())
	ret = append(ret, IndexUpdateAllErrorCnt.Get())

	// index update incr
	ret = append(ret, IndexUpdateIncr.Get())
	ret = append(ret, IndexUpdateIncrCnt.Get())
	ret = append(ret, IndexUpdateIncrErrorCnt.Get())

	ret = append(ret, IndexUpdateIncrDbEndpointSelectCnt.Get())
	ret = append(ret, IndexUpdateIncrDbEndpointInsertCnt.Get())

	ret = append(ret, IndexUpdateIncrDbTagEndpointSelectCnt.Get())
	ret = append(ret, IndexUpdateIncrDbTagEndpointInsertCnt.Get())

	ret = append(ret, IndexUpdateIncrDbEndpointCounterSelectCnt.Get())
	ret = append(ret, IndexUpdateIncrDbEndpointCounterInsertCnt.Get())
	ret = append(ret, IndexUpdateIncrDbEndpointCounterUpdateCnt.Get())

	// index db cache
	ret = append(ret, IndexedItemCacheCnt.Get())
	ret = append(ret, UnIndexedItemCacheCnt.Get())
	ret = append(ret, IndexDbEndpointCacheCnt.Get())
	ret = append(ret, IndexDbTagEndpointCacheCnt.Get())
	ret = append(ret, IndexDbEndpointCounterCacheCnt.Get())

	return ret
}

// TODO 临时放在这里了, 考虑放到合适的模块
func FmtUnixTs(ts int64) string {
	return time.Unix(ts, 0).Format("2006-01-02 15:04:05")
}
