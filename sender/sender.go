package sender

import (
	"fmt"
	cmodel "github.com/open-falcon/common/model"
	"github.com/open-falcon/transfer/g"
	"github.com/open-falcon/transfer/proc"
	cpool "github.com/open-falcon/transfer/sender/conn_pool"
	"github.com/toolkits/container/list"
	"log"
)

const (
	DefaultSendQueueMaxSize = 10240 //1.024w
)

// 服务节点的一致性哈希环
// pk -> node
var (
	JudgeNodeRing          *ConsistentHashNodeRing
	GraphNodeRing          *ConsistentHashNodeRing
	GraphMigratingNodeRing *ConsistentHashNodeRing
)

// 发送缓存队列
// node -> queue_of_data
var (
	JudgeQueues          = make(map[string]*list.SafeLinkedListLimited)
	GraphQueues          = make(map[string]*list.SafeLinkedListLimited)
	GraphMigratingQueues = make(map[string]*list.SafeLinkedListLimited)
)

// 连接池
// node_address -> connection_pool
var (
	JudgeConnPools          *cpool.SafeRpcConnPools
	GraphConnPools          *cpool.SafeRpcConnPools
	GraphMigratingConnPools *cpool.SafeRpcConnPools
)

// 初始化数据发送服务, 在main函数中调用
func Start() {
	initConnPools()
	initSendQueues()
	initNodeRings()
	// SendTasks依赖基础组件的初始化,要最后启动
	startSendTasks()
	startSenderCron()
	log.Println("send.Start, ok")
}

// 将数据 打入 某个Judge的发送缓存队列, 具体是哪一个Judge 由一致性哈希 决定
func Push2JudgeSendQueue(items []*cmodel.MetaData) {
	for _, item := range items {
		pk := item.PK()

		// statistics, transfer recv. 为了效率,放到了这里,因此只有judge是enbale时才能trace
		proc.RecvDataTrace.Trace(pk, item)

		node, err := JudgeNodeRing.GetNode(pk)
		if err != nil {
			log.Println("E:", err)
			continue
		}

		judgeItem := &cmodel.JudgeItem{
			Endpoint:  item.Endpoint,
			Metric:    item.Metric,
			Value:     item.Value,
			Timestamp: item.Timestamp,
			JudgeType: item.CounterType,
			Tags:      item.Tags,
		}
		Q := JudgeQueues[node]
		isSuccess := Q.PushFront(judgeItem)

		// statistics
		if !isSuccess {
			proc.SendToJudgeDropCnt.Incr()
		}
	}
}

// 将数据 打入 某个Graph的发送缓存队列, 具体是哪一个Graph 由一致性哈希 决定
// 如果正在数据迁移, 数据除了打到原有配置上 还要向新的配置上打一份(新老重叠时要去重,防止将一条数据向一台Graph上打两次)
func Push2GraphSendQueue(items []*cmodel.MetaData, migrating bool) {
	for _, item := range items {
		graphItem, err := convert2GraphItem(item)
		if err != nil {
			log.Println("E:", err)
			continue
		}

		pk := item.PK()
		node, err := GraphNodeRing.GetNode(pk)
		if err != nil {
			log.Println("E:", err)
			continue
		}
		Q := GraphQueues[node]
		isSuccess := Q.PushFront(graphItem)

		// statistics
		if !isSuccess {
			proc.SendToGraphDropCnt.Incr()
		}

		if migrating {
			migratingNode, err := GraphMigratingNodeRing.GetNode(pk)
			if err != nil {
				log.Println("E:", err)
				continue
			}

			if node != migratingNode { // 数据迁移时,进行新老去重
				MQ := GraphMigratingQueues[migratingNode]
				isSuccess := MQ.PushFront(graphItem)

				// statistics
				if !isSuccess {
					proc.SendToGraphMigratingDropCnt.Incr()

				}
			}
		}
	}
}

// 打到Graph的数据,要根据rrdtool的特定 来限制 step、counterType、timestamp
func convert2GraphItem(d *cmodel.MetaData) (*cmodel.GraphItem, error) {
	item := &cmodel.GraphItem{}

	item.Endpoint = d.Endpoint
	item.Metric = d.Metric
	item.Tags = d.Tags
	item.Timestamp = d.Timestamp
	item.Value = d.Value
	item.Step = int(d.Step)
	if item.Step < g.MIN_STEP {
		item.Step = g.MIN_STEP
	}
	item.Heartbeat = item.Step * 2

	if d.CounterType == g.GAUGE {
		item.DsType = d.CounterType
		item.Min = "U"
		item.Max = "U"
	} else if d.CounterType == g.COUNTER {
		item.DsType = g.DERIVE
		item.Min = "0"
		item.Max = "U"
	} else if d.CounterType == g.DERIVE {
		item.DsType = g.DERIVE
		item.Min = "0"
		item.Max = "U"
	} else {
		return item, fmt.Errorf("not_supported_counter_type")
	}

	item.Timestamp = item.Timestamp - item.Timestamp%int64(item.Step)

	return item, nil
}
