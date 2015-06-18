package proc

import (
	nproc "github.com/toolkits/proc"
)

// trace
var (
	RecvDataTrace = nproc.NewDataTrace("RecvDataTrace", 3)
)

// filter
var (
	RecvDataFilter = nproc.NewDataFilter("RecvDataFilter", 3)
)

// 索引增量更新
var (
	IndexUpdateIncr         = nproc.NewSCounterQps("IndexUpdateIncr")
	IndexUpdateIncrCnt      = nproc.NewSCounterBase("IndexUpdateIncrCnt")
	IndexUpdateIncrErrorCnt = nproc.NewSCounterQps("IndexUpdateIncrErrorCnt")

	IndexUpdateIncrDbEndpointSelectCnt = nproc.NewSCounterQps("IndexUpdateIncrDbEndpointSelectCnt")
	IndexUpdateIncrDbEndpointInsertCnt = nproc.NewSCounterQps("IndexUpdateIncrDbEndpointInsertCnt")

	IndexUpdateIncrDbTagEndpointSelectCnt = nproc.NewSCounterQps("IndexUpdateIncrDbTagEndpointSelectCnt")
	IndexUpdateIncrDbTagEndpointInsertCnt = nproc.NewSCounterQps("IndexUpdateIncrDbTagEndpointInsertCnt")

	IndexUpdateIncrDbEndpointCounterSelectCnt = nproc.NewSCounterQps("IndexUpdateIncrDbEndpointCounterSelectCnt")
	IndexUpdateIncrDbEndpointCounterInsertCnt = nproc.NewSCounterQps("IndexUpdateIncrDbEndpointCounterInsertCnt")
	IndexUpdateIncrDbEndpointCounterUpdateCnt = nproc.NewSCounterQps("IndexUpdateIncrDbEndpointCounterUpdateCnt")
)

// 索引全量更新
var (
	IndexUpdateAll         = nproc.NewSCounterQps("IndexUpdateAll")
	IndexUpdateAllCnt      = nproc.NewSCounterBase("IndexUpdateAllCnt")
	IndexUpdateAllErrorCnt = nproc.NewSCounterQps("IndexUpdateAllErrorCnt")
)

// 索引缓存大小
var (
	IndexedItemCacheCnt   = nproc.NewSCounterBase("IndexedItemCacheCnt")
	UnIndexedItemCacheCnt = nproc.NewSCounterBase("UnIndexedItemCacheCnt")
	EndpointCacheCnt      = nproc.NewSCounterBase("EndpointCacheCnt")
	CounterCacheCnt       = nproc.NewSCounterBase("CounterCacheCnt")
)

// Rpc
var (
	GraphRpcRecvCnt = nproc.NewSCounterQps("GraphRpcRecvCnt")
)

// Query
var (
	GraphQueryCnt     = nproc.NewSCounterQps("GraphQueryCnt")
	GraphQueryItemCnt = nproc.NewSCounterQps("GraphQueryItemCnt")
	GraphInfoCnt      = nproc.NewSCounterQps("GraphInfoCnt")
	GraphLastCnt      = nproc.NewSCounterQps("GraphLastCnt")
	GraphLoadDbCnt    = nproc.NewSCounterQps("GraphLoadDbCnt") // load sth from db when query/info, tmp
)

func GetAll() []interface{} {
	ret := make([]interface{}, 0)

	// rpc recv
	ret = append(ret, GraphRpcRecvCnt.Get())

	// query
	ret = append(ret, GraphQueryCnt.Get())
	ret = append(ret, GraphQueryItemCnt.Get())
	ret = append(ret, GraphInfoCnt.Get())
	ret = append(ret, GraphLastCnt.Get())
	ret = append(ret, GraphLoadDbCnt.Get())

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
	ret = append(ret, EndpointCacheCnt.Get())
	ret = append(ret, CounterCacheCnt.Get())

	return ret
}
