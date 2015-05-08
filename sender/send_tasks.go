package sender

import (
	"github.com/open-falcon/common/model"
	"github.com/open-falcon/transfer/g"
	"github.com/open-falcon/transfer/proc"
	tSema "github.com/toolkits/concurrent/semaphore"
	"github.com/toolkits/container/list"
	"log"
	"time"
)

const (
	DefaultSendTaskSleepInterval = time.Millisecond * 50 //默认睡眠间隔为50ms
)

var (
	semaSendToJudge, semaSendToGraph, semaSendToGraphMigrating *tSema.Semaphore
)

// TODO 添加对发送任务的控制,比如stop等
func startSendTasks() {
	cfg := g.Config()
	// init semaphore
	judgeConcurrent := cfg.Judge.MaxIdle / 2
	graphConcurrent := cfg.Graph.MaxIdle / 2
	if judgeConcurrent < 1 {
		judgeConcurrent = 1
	}
	if graphConcurrent < 1 {
		graphConcurrent = 1
	}
	semaSendToJudge = tSema.NewSemaphore(judgeConcurrent)
	semaSendToGraph = tSema.NewSemaphore(graphConcurrent)
	semaSendToGraphMigrating = tSema.NewSemaphore(graphConcurrent)

	// tasks
	for node, _ := range cfg.Judge.Cluster {
		queue := JudgeQueues[node]
		go forward2JudgeTask(queue, node)
	}

	for node, _ := range cfg.Graph.Cluster {
		queue := GraphQueues[node]
		go forward2GraphTask(queue, node)
	}

	if cfg.Graph.Migrating {
		for node, _ := range cfg.Graph.ClusterMigrating {
			queue := GraphMigratingQueues[node]
			go forward2GraphMigratingTask(queue, node)
		}
	}
}

// Judge定时任务, 将 Judge发送缓存中的数据 通过rpc连接池 发送到Judge
func forward2JudgeTask(Q *list.SafeLinkedListLimited, node string) {
	batch := g.Config().Judge.Batch // 一次发送,最多batch条数据
	addr := g.Config().Judge.Cluster[node]

	for {
		items := Q.PopBack(batch)
		count := len(items)
		if count == 0 {
			time.Sleep(DefaultSendTaskSleepInterval)
			continue
		}

		judgeItems := make([]*model.JudgeItem, count)
		for i := 0; i < count; i++ {
			judgeItems[i] = items[i].(*model.JudgeItem)
		}

		//	同步Call + 有限并发 进行发送
		semaSendToJudge.Acquire()
		go func(addr string, judgeItems []*model.JudgeItem, count int) {
			resp := &model.SimpleRpcResponse{}
			err := JudgeConnPools.Call(addr, "Judge.Send", judgeItems, resp)
			if err != nil {
				log.Printf("send judge %s fail: %v", addr, err)
				// statistics
				proc.SendToJudgeFailCnt.IncrBy(int64(count))
			} else {
				// statistics
				proc.SendToJudgeCnt.IncrBy(int64(count))
				cnt := proc.SendToJudgeCntPerNode[node]
				if cnt != nil {
					cnt.IncrBy(int64(count))
				}
			}
			semaSendToJudge.Release()
		}(addr, judgeItems, count)
	}
}

// Graph定时任务, 将 Graph发送缓存中的数据 通过rpc连接池 发送到Graph
func forward2GraphTask(Q *list.SafeLinkedListLimited, node string) {
	batch := g.Config().Graph.Batch // 一次发送,最多batch条数据
	addr := g.Config().Graph.Cluster[node]

	for {
		items := Q.PopBack(batch)
		count := len(items)
		if count == 0 {
			time.Sleep(DefaultSendTaskSleepInterval)
			continue
		}

		graphItems := make([]*model.GraphItem, count)
		for i := 0; i < count; i++ {
			graphItems[i] = items[i].(*model.GraphItem)
		}

		semaSendToGraph.Acquire()
		go func(addr string, graphItems []*model.GraphItem, count int) {
			resp := &model.SimpleRpcResponse{}
			err := GraphConnPools.Call(addr, "Graph.Send", graphItems, resp)
			if err != nil {
				log.Printf("send to graph %s fail: %v", addr, err)
				// statistics
				proc.SendToGraphFailCnt.IncrBy(int64(count))
			} else {
				// statistics
				proc.SendToGraphCnt.IncrBy(int64(count))
				cnt := proc.SendToGraphCntPerNode[node]
				if cnt != nil {
					cnt.IncrBy(int64(count))
				}
			}
			semaSendToGraph.Release()
		}(addr, graphItems, count)
	}
}

// Graph定时任务, 进行数据迁移时的 数据冗余发送
func forward2GraphMigratingTask(Q *list.SafeLinkedListLimited, node string) {
	batch := g.Config().Graph.Batch // 一次发送,最多batch条数据
	addr := g.Config().Graph.ClusterMigrating[node]

	for {
		items := Q.PopBack(batch)
		count := len(items)
		if count == 0 {
			time.Sleep(DefaultSendTaskSleepInterval)
			continue
		}

		graphItems := make([]*model.GraphItem, count)
		for i := 0; i < count; i++ {
			graphItems[i] = items[i].(*model.GraphItem)
		}

		semaSendToGraphMigrating.Acquire()
		go func(addr string, graphItems []*model.GraphItem, count int) {
			resp := &model.SimpleRpcResponse{}
			err := GraphMigratingConnPools.Call(addr, "Graph.Send", graphItems, resp)
			if err != nil {
				log.Printf("send to graph migrating %s fail: %v", addr, err)
				// statistics
				proc.SendToGraphMigratingFailCnt.IncrBy(int64(count))
			} else {
				// statistics
				proc.SendToGraphMigratingCnt.IncrBy(int64(count))
				cnt := proc.SendToGraphMigratingCntPerNode[node]
				if cnt != nil {
					cnt.IncrBy(int64(count))
				}
			}
			semaSendToGraphMigrating.Release()
		}(addr, graphItems, count)
	}
}
