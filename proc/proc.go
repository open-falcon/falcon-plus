package proc

import (
	MP "github.com/open-falcon/common/proc"
	"github.com/open-falcon/transfer/g"
	"log"
)

const (
	DefaultProcCacheSize     = 3 // 写死, 用于自运维的 数据缓存的 大小,如发送缓存、接收缓存等
	DefaultSCounterQpsPeriod = 5 // QPS计算周期, 默认值为1s
)

// 已接收数据的缓存,用于调试. 只缓存一条监控数据的N个点

// 统计指标的整体数据
var (
	// 计数统计,正确计数,错误计数, ...
	RecvCnt       = MP.NewSCounterQps("RecvCnt")
	RpcRecvCnt    = MP.NewSCounterQps("RpcRecvCnt")
	SocketRecvCnt = MP.NewSCounterQps("SocketRecvCnt")

	SendToJudgeCnt          = MP.NewSCounterQps("SendToJudgeCnt")
	SendToGraphCnt          = MP.NewSCounterQps("SendToGraphCnt")
	SendToGraphMigratingCnt = MP.NewSCounterQps("SendToGraphMigratingCnt")

	SendToJudgeDropCnt          = MP.NewSCounterQps("SendToJudgeDropCnt")
	SendToGraphDropCnt          = MP.NewSCounterQps("SendToGraphDropCnt")
	SendToGraphMigratingDropCnt = MP.NewSCounterQps("SendToGraphMigratingDropCnt")

	SendToJudgeFailCnt          = MP.NewSCounterQps("SendToJudgeFailCnt")
	SendToGraphFailCnt          = MP.NewSCounterQps("SendToGraphFailCnt")
	SendToGraphMigratingFailCnt = MP.NewSCounterQps("SendToGraphMigratingFailCnt")

	// 发送缓存大小
	JudgeQueuesCnt          = MP.NewSCounterBase("JudgeSendCacheCnt")
	GraphQueuesCnt          = MP.NewSCounterBase("GraphSendCacheCnt")
	GraphMigratingQueuesCnt = MP.NewSCounterBase("GraphMigratingCacheCnt")
)

var (
	// 发送到每个节点的 计数统计
	SendToJudgeCntPerNode          map[string]*MP.SCounterQps
	SendToGraphCntPerNode          map[string]*MP.SCounterQps
	SendToGraphMigratingCntPerNode map[string]*MP.SCounterQps

	// 发送到每个节点的丢失情况的 计数统计
	SendToJudgeDropCntPerNode          map[string]*MP.SCounterQps
	SendToGraphDropCntPerNode          map[string]*MP.SCounterQps
	SendToGraphMigratingDropCntPerNode map[string]*MP.SCounterQps

	// 每个节点的发送缓存大小
	JudgeQueuesCntPerNode          map[string]*MP.SCounterBase
	GraphQueuesCntPerNode          map[string]*MP.SCounterBase
	GraphMigratingQueuesCntPerNode map[string]*MP.SCounterBase
)

func Init() {
	// proc存储容器初始化
	cfg := g.Config()
	judgeCluster := cfg.Judge.Cluster
	graphCluster := cfg.Graph.Cluster
	graphMigratingCluster := cfg.Graph.ClusterMigrating

	SendToJudgeCntPerNode = make(map[string]*MP.SCounterQps)
	SendToJudgeDropCntPerNode = make(map[string]*MP.SCounterQps)
	JudgeQueuesCntPerNode = make(map[string]*MP.SCounterBase)
	for node, _ := range judgeCluster {
		cnt := MP.NewSCounterQps("SendToJudgeCntPerNode." + node)
		SendToJudgeCntPerNode[node] = cnt

		dropCnt := MP.NewSCounterQps("SendToJudgeDropCntPerNode." + node)
		SendToJudgeDropCntPerNode[node] = dropCnt

		queueSize := MP.NewSCounterBase("JudgeQueuesCntPerNode." + node)
		JudgeQueuesCntPerNode[node] = queueSize
	}

	SendToGraphCntPerNode = make(map[string]*MP.SCounterQps)
	SendToGraphDropCntPerNode = make(map[string]*MP.SCounterQps)
	GraphQueuesCntPerNode = make(map[string]*MP.SCounterBase)
	for node, _ := range graphCluster {
		cnt := MP.NewSCounterQps("SendToGraphCntPerNode." + node)
		SendToGraphCntPerNode[node] = cnt

		dropCnt := MP.NewSCounterQps("SendToGraphDropCntPerNode." + node)
		SendToGraphDropCntPerNode[node] = dropCnt

		queueSize := MP.NewSCounterBase("GraphQueuesCntPerNode." + node)
		GraphQueuesCntPerNode[node] = queueSize
	}

	SendToGraphMigratingCntPerNode = make(map[string]*MP.SCounterQps)
	SendToGraphMigratingDropCntPerNode = make(map[string]*MP.SCounterQps)
	GraphMigratingQueuesCntPerNode = make(map[string]*MP.SCounterBase)
	for node, _ := range graphMigratingCluster {
		cnt := MP.NewSCounterQps("SendToGraphMigratingCntPerNode." + node)
		SendToGraphMigratingCntPerNode[node] = cnt

		dropCnt := MP.NewSCounterQps("SendToGraphMigratingDropCntPerNode." + node)
		SendToGraphMigratingDropCntPerNode[node] = dropCnt

		queueSize := MP.NewSCounterBase("GraphMigratingQueuesCntPerNode." + node)
		GraphMigratingQueuesCntPerNode[node] = queueSize
	}

	log.Println("proc.Init, ok")
}

func GetAll() []interface{} {
	ret := make([]interface{}, 0)

	// recv cnt
	ret = append(ret, RecvCnt)
	ret = append(ret, RpcRecvCnt)
	ret = append(ret, SocketRecvCnt)

	// send cnt
	ret = append(ret, SendToJudgeCnt)
	ret = append(ret, SendToGraphCnt)
	ret = append(ret, SendToGraphMigratingCnt)

	// drop cnt
	ret = append(ret, SendToJudgeDropCnt)
	ret = append(ret, SendToGraphDropCnt)
	ret = append(ret, SendToGraphMigratingDropCnt)

	// send fail cnt
	ret = append(ret, SendToJudgeFailCnt)
	ret = append(ret, SendToGraphFailCnt)
	ret = append(ret, SendToGraphMigratingFailCnt)

	// cache cnt
	ret = append(ret, JudgeQueuesCnt)
	ret = append(ret, GraphQueuesCnt)
	ret = append(ret, GraphMigratingQueuesCnt)

	return ret
}
