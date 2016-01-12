package sender

import (
	"fmt"
	cmodel "github.com/open-falcon/common/model"
	"github.com/open-falcon/transfer/g"
	"github.com/open-falcon/transfer/proc"
	cpool "github.com/open-falcon/transfer/sender/conn_pool"
	nlist "github.com/toolkits/container/list"
	"log"
)

const (
	DefaultSendQueueMaxSize = 102400 //10.24w
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
	TsdbQueue            *nlist.SafeListLimited
	JudgeQueues          = make(map[string]*nlist.SafeListLimited)
	GraphQueues          = make(map[string]*nlist.SafeListLimited)
	GraphMigratingQueues = make(map[string]*nlist.SafeListLimited)
)

// 连接池
// node_address -> connection_pool
var (
	JudgeConnPools          *cpool.SafeRpcConnPools
	TsdbConnPools           *cpool.SafeTcpConnPools
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
		node, err := JudgeNodeRing.GetNode(pk)
		if err != nil {
			log.Println("E:", err)
			continue
		}

		// align ts
		step := int(item.Step)
		if step < g.MIN_STEP {
			step = g.MIN_STEP
		}
		ts := alignTs(item.Timestamp, int64(step))

		judgeItem := &cmodel.JudgeItem{
			Endpoint:  item.Endpoint,
			Metric:    item.Metric,
			Value:     item.Value,
			Timestamp: ts,
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
	cfg := g.Config().Graph

	for _, item := range items {
		graphItem, err := convert2GraphItem(item)
		if err != nil {
			log.Println("E:", err)
			continue
		}
		pk := item.PK()

		// statistics. 为了效率,放到了这里,因此只有graph是enbale时才能trace
		proc.RecvDataTrace.Trace(pk, item)
		proc.RecvDataFilter.Filter(pk, item.Value, item)

		node, err := GraphNodeRing.GetNode(pk)
		if err != nil {
			log.Println("E:", err)
			continue
		}

		cnode := cfg.Cluster2[node]
		errCnt := 0
		for _, addr := range cnode.Addrs {
			Q := GraphQueues[node+addr]
			if !Q.PushFront(graphItem) {
				errCnt += 1
			}
		}

		// statistics
		if errCnt > 0 {
			proc.SendToGraphDropCnt.Incr()
		}

		if migrating {
			migratingNode, err := GraphMigratingNodeRing.GetNode(pk)
			if err != nil {
				log.Println("E:", err)
				continue
			}

			if node != migratingNode { // 数据迁移时,进行新老去重
				cnodem := cfg.ClusterMigrating2[migratingNode]
				errCnt := 0
				for _, addr := range cnodem.Addrs {
					MQ := GraphMigratingQueues[migratingNode+addr]
					if !MQ.PushFront(graphItem) {
						errCnt += 1
					}
				}

				// statistics
				if errCnt > 0 {
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

	item.Timestamp = alignTs(item.Timestamp, int64(item.Step)) //item.Timestamp - item.Timestamp%int64(item.Step)

	return item, nil
}

// 将原始数据入到tsdb发送缓存队列
func Push2TsdbSendQueue(items []*cmodel.MetaData) {
	for _, item := range items {
		tsdbItem := convert2TsdbItem(item)
		isSuccess := TsdbQueue.PushFront(tsdbItem)

		if !isSuccess {
			proc.SendToTsdbDropCnt.Incr()
		}
	}
}

// 转化为tsdb格式
func convert2TsdbItem(d *cmodel.MetaData) *cmodel.TsdbItem {
	t := cmodel.TsdbItem{Tags: make(map[string]string)}

	for k, v := range d.Tags {
		t.Tags[k] = v
	}
	t.Tags["endpoint"] = d.Endpoint
	t.Metric = d.Metric
	t.Timestamp = d.Timestamp
	t.Value = d.Value
	return &t
}

func alignTs(ts int64, period int64) int64 {
	return ts - ts%period
}
