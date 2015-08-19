package proc

import (
	"log"

	nproc "github.com/toolkits/proc"
)

// trace
var (
	RecvDataTrace = nproc.NewDataTrace("RecvDataTrace", 5)
)

// filter
var (
	RecvDataFilter = nproc.NewDataFilter("RecvDataFilter", 5)
)

// counter
var (
	// recv
	RecvCnt       = nproc.NewSCounterQps("RecvCnt")
	RpcRecvCnt    = nproc.NewSCounterQps("RpcRecvCnt")
	HttpRecvCnt   = nproc.NewSCounterQps("HttpRecvCnt")
	SocketRecvCnt = nproc.NewSCounterQps("SocketRecvCnt")

	// send
	SendCnt     = nproc.NewSCounterQps("SendCnt")
	SendDropCnt = nproc.NewSCounterQps("SendDropCnt")
	SendFailCnt = nproc.NewSCounterQps("SendFailCnt")

	SendQueuesCnt = nproc.NewSCounterBase("SendQueuesCnt")
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
	ret = append(ret, SendCnt.Get())
	ret = append(ret, SendDropCnt.Get())
	ret = append(ret, SendFailCnt.Get())

	ret = append(ret, SendQueuesCnt.Get())

	return ret
}
