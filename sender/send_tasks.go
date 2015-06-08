package sender

import (
	nsema "github.com/niean/gotools/concurrent/semaphore"
	cmodel "github.com/open-falcon/common/model"
	"github.com/open-falcon/transfer/g"
	"github.com/open-falcon/transfer/proc"
	"github.com/toolkits/container/list"
	"log"
	"time"
)

// send
const (
	DefaultSendTaskSleepInterval = time.Millisecond * 50 //默认睡眠间隔为50ms
)

var (
	semaSendToJudge, semaSendToGraph, semaSendToGraphMigrating *nsema.Semaphore
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
	semaSendToJudge = nsema.NewSemaphore(judgeConcurrent)
	semaSendToGraph = nsema.NewSemaphore(graphConcurrent)
	semaSendToGraphMigrating = nsema.NewSemaphore(graphConcurrent)

	// init send go-routines
	for node, _ := range cfg.Judge.Cluster {
		queue := JudgeQueues[node]
		go forward2JudgeTask(queue, node, judgeConcurrent)
	}

	for node, _ := range cfg.Graph.Cluster {
		queue := GraphQueues[node]
		go forward2GraphTask(queue, node, graphConcurrent)
	}

	if cfg.Graph.Migrating {
		for node, _ := range cfg.Graph.ClusterMigrating {
			queue := GraphMigratingQueues[node]
			go forward2GraphMigratingTask(queue, node, graphConcurrent)
		}
	}
}

// Judge定时任务, 将 Judge发送缓存中的数据 通过rpc连接池 发送到Judge
func forward2JudgeTask(Q *list.SafeLinkedListLimited, node string, concurrent int) {
	batch := g.Config().Judge.Batch // 一次发送,最多batch条数据
	addr := g.Config().Judge.Cluster[node]
	sema := nsema.NewSemaphore(concurrent)

	for {
		items := Q.PopBack(batch)
		count := len(items)
		if count == 0 {
			time.Sleep(DefaultSendTaskSleepInterval)
			continue
		}

		judgeItems := make([]*cmodel.JudgeItem, count)
		for i := 0; i < count; i++ {
			judgeItems[i] = items[i].(*cmodel.JudgeItem)
		}

		//	同步Call + 有限并发 进行发送
		sema.Acquire()
		go func(addr string, judgeItems []*cmodel.JudgeItem, count int) {
			defer sema.Release()

			resp := &cmodel.SimpleRpcResponse{}
			var err error
			sendOk := false
			for i := 0; i < 3; i++ { //最多重试3次
				err = JudgeConnPools.Call(addr, "Judge.Send", judgeItems, resp)
				if err == nil {
					sendOk = true
					break
				}
				time.Sleep(time.Millisecond * 10)
			}

			if !sendOk {
				log.Printf("send judge %s fail: %v", addr, err)
				// statistics
				proc.SendToJudgeFailCnt.IncrBy(int64(count))
			} else {
				// statistics
				proc.SendToJudgeCnt.IncrBy(int64(count))
			}
		}(addr, judgeItems, count)
	}
}

// Graph定时任务, 将 Graph发送缓存中的数据 通过rpc连接池 发送到Graph
func forward2GraphTask(Q *list.SafeLinkedListLimited, node string, concurrent int) {
	batch := g.Config().Graph.Batch // 一次发送,最多batch条数据
	addr := g.Config().Graph.Cluster[node]
	sema := nsema.NewSemaphore(concurrent)

	for {
		items := Q.PopBack(batch)
		count := len(items)
		if count == 0 {
			time.Sleep(DefaultSendTaskSleepInterval)
			continue
		}

		graphItems := make([]*cmodel.GraphItem, count)
		for i := 0; i < count; i++ {
			graphItems[i] = items[i].(*cmodel.GraphItem)
		}

		sema.Acquire()
		go func(addr string, graphItems []*cmodel.GraphItem, count int) {
			defer sema.Release()

			resp := &cmodel.SimpleRpcResponse{}
			var err error
			sendOk := false
			for i := 0; i < 3; i++ { //最多重试3次
				err = GraphConnPools.Call(addr, "Graph.Send", graphItems, resp)
				if err == nil {
					sendOk = true
					break
				}
				time.Sleep(time.Millisecond * 10)
			}

			if !sendOk {
				log.Printf("send to graph %s fail: %v", addr, err)
				// statistics
				proc.SendToGraphFailCnt.IncrBy(int64(count))
			} else {
				// statistics
				proc.SendToGraphCnt.IncrBy(int64(count))
			}
		}(addr, graphItems, count)
	}
}

// Graph定时任务, 进行数据迁移时的 数据冗余发送
func forward2GraphMigratingTask(Q *list.SafeLinkedListLimited, node string, concurrent int) {
	batch := g.Config().Graph.Batch // 一次发送,最多batch条数据
	addr := g.Config().Graph.ClusterMigrating[node]
	sema := nsema.NewSemaphore(concurrent)

	for {
		items := Q.PopBack(batch)
		count := len(items)
		if count == 0 {
			time.Sleep(DefaultSendTaskSleepInterval)
			continue
		}

		graphItems := make([]*cmodel.GraphItem, count)
		for i := 0; i < count; i++ {
			graphItems[i] = items[i].(*cmodel.GraphItem)
		}

		sema.Acquire()
		go func(addr string, graphItems []*cmodel.GraphItem, count int) {
			defer sema.Release()

			resp := &cmodel.SimpleRpcResponse{}
			var err error
			sendOk := false
			for i := 0; i < 3; i++ { //最多重试3次
				err = GraphMigratingConnPools.Call(addr, "Graph.Send", graphItems, resp)
				if err == nil {
					sendOk = true
					break
				}
				time.Sleep(time.Millisecond * 10) //发送失败了,睡10ms
			}

			if !sendOk {
				log.Printf("send to graph migrating %s fail: %v", addr, err)
				// statistics
				proc.SendToGraphMigratingFailCnt.IncrBy(int64(count))
			} else {
				// statistics
				proc.SendToGraphMigratingCnt.IncrBy(int64(count))
			}
		}(addr, graphItems, count)
	}
}
