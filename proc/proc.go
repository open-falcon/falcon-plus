package proc

import (
	nproc "github.com/toolkits/proc"
	"log"
)

// trace
var (
	RecvDataTrace = nproc.NewDataTrace("RecvDataTrace", 3)
)

// filter
var (
	RecvDataFilter = nproc.NewDataFilter("RecvDataFilter", 5)
)

// 统计指标的整体数据
var (
	// 计数统计,正确计数,错误计数, ...
	RecvCnt       = nproc.NewSCounterQps("RecvCnt")
	RpcRecvCnt    = nproc.NewSCounterQps("RpcRecvCnt")
	HttpRecvCnt   = nproc.NewSCounterQps("HttpRecvCnt")
	SocketRecvCnt = nproc.NewSCounterQps("SocketRecvCnt")

	SendToJudgeCnt = nproc.NewSCounterQps("SendToJudgeCnt")
	SendToTsdbCnt  = nproc.NewSCounterQps("SendToTsdbCnt")
	SendToGraphCnt = nproc.NewSCounterQps("SendToGraphCnt")

	SendToJudgeDropCnt = nproc.NewSCounterQps("SendToJudgeDropCnt")
	SendToTsdbDropCnt  = nproc.NewSCounterQps("SendToTsdbDropCnt")
	SendToGraphDropCnt = nproc.NewSCounterQps("SendToGraphDropCnt")

	SendToJudgeFailCnt = nproc.NewSCounterQps("SendToJudgeFailCnt")
	SendToTsdbFailCnt  = nproc.NewSCounterQps("SendToTsdbFailCnt")
	SendToGraphFailCnt = nproc.NewSCounterQps("SendToGraphFailCnt")

	// 发送缓存大小
	JudgeQueuesCnt = nproc.NewSCounterBase("JudgeSendCacheCnt")
	TsdbQueuesCnt  = nproc.NewSCounterBase("TsdbSendCacheCnt")
	GraphQueuesCnt = nproc.NewSCounterBase("GraphSendCacheCnt")
)

func Start() {
	log.Println("proc.Start, ok")
}

func GetAll() []interface{} {
	ret := make([]interface{}, 0)

	// recv cnt
	ret = append(ret, RecvCnt.Get())
	ret = append(ret, RpcRecvCnt.Get())
	ret = append(ret, HttpRecvCnt.Get())
	ret = append(ret, SocketRecvCnt.Get())

	// send cnt
	ret = append(ret, SendToJudgeCnt.Get())
	ret = append(ret, SendToTsdbCnt.Get())
	ret = append(ret, SendToGraphCnt.Get())

	// drop cnt
	ret = append(ret, SendToJudgeDropCnt.Get())
	ret = append(ret, SendToTsdbDropCnt.Get())
	ret = append(ret, SendToGraphDropCnt.Get())

	// send fail cnt
	ret = append(ret, SendToJudgeFailCnt.Get())
	ret = append(ret, SendToTsdbFailCnt.Get())
	ret = append(ret, SendToGraphFailCnt.Get())

	// cache cnt
	ret = append(ret, JudgeQueuesCnt.Get())
	ret = append(ret, TsdbQueuesCnt.Get())
	ret = append(ret, GraphQueuesCnt.Get())

	return ret
}
