// Copyright 2017 Xiaomi, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sender

import (
	"bytes"
	"log"
	"math/rand"
	"time"

	pfc "github.com/niean/goperfcounter"

	cmodel "github.com/open-falcon/falcon-plus/common/model"
	cutils "github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/transfer/g"
	"github.com/open-falcon/falcon-plus/modules/transfer/proc"
	nsema "github.com/toolkits/concurrent/semaphore"
	"github.com/toolkits/container/list"
)

// send
const (
	DefaultSendTaskSleepInterval = time.Millisecond * 50 //默认睡眠间隔为50ms
)

// TODO 添加对发送任务的控制,比如stop等
func startSendTasks() {
	cfg := g.Config()
	// init semaphore
	judgeConcurrent := cfg.Judge.MaxConns
	graphConcurrent := cfg.Graph.MaxConns
	tsdbConcurrent := cfg.Tsdb.MaxConns
	transferConcurrent := cfg.Transfer.MaxConns

	if tsdbConcurrent < 1 {
		tsdbConcurrent = 1
	}

	if judgeConcurrent < 1 {
		judgeConcurrent = 1
	}

	if graphConcurrent < 1 {
		graphConcurrent = 1
	}

	if transferConcurrent < 1 {
		transferConcurrent = 1
	}

	if cfg.Judge.Enabled {
		// init send go-routines
		for node := range cfg.Judge.Cluster {
			queue := JudgeQueues[node]
			go forward2JudgeTask(queue, node, judgeConcurrent)
		}
	}

	if cfg.Graph.Enabled {
		for node, nitem := range cfg.Graph.ClusterList {
			for _, addr := range nitem.Addrs {
				queue := GraphQueues[node+addr]
				go forward2GraphTask(queue, node, addr, graphConcurrent)
			}
		}
	}

	if cfg.Transfer.Enabled {
		// init transfer send cnt
		concurrent := cfg.Transfer.MaxConns * len(cfg.Transfer.Cluster)
		go forward2TransferTask(TransferQueue, concurrent)
	}

	if cfg.Tsdb.Enabled {
		go forward2TsdbTask(tsdbConcurrent)
	}
}

// Judge定时任务, 将 Judge发送缓存中的数据 通过rpc连接池 发送到Judge
func forward2JudgeTask(Q *list.SafeListLimited, node string, concurrent int) {
	batch := g.Config().Judge.Batch // 一次发送,最多batch条数据
	addr := g.Config().Judge.Cluster[node]
	sema := nsema.NewSemaphore(concurrent)

	for {
		items := Q.PopBackBy(batch)
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

			// statistics
			if !sendOk {
				log.Printf("send judge %s:%s fail: %v", node, addr, err)
				proc.SendToJudgeFailCnt.IncrBy(int64(count))
			} else {
				proc.SendToJudgeCnt.IncrBy(int64(count))
			}
		}(addr, judgeItems, count)
	}
}

// Graph定时任务, 将 Graph发送缓存中的数据 通过rpc连接池 发送到Graph
func forward2GraphTask(Q *list.SafeListLimited, node string, addr string, concurrent int) {
	batch := g.Config().Graph.Batch // 一次发送,最多batch条数据
	sema := nsema.NewSemaphore(concurrent)

	for {
		items := Q.PopBackBy(batch)
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

			// statistics
			if !sendOk {
				log.Printf("send to graph %s:%s fail: %v", node, addr, err)
				proc.SendToGraphFailCnt.IncrBy(int64(count))
			} else {
				proc.SendToGraphCnt.IncrBy(int64(count))
			}
		}(addr, graphItems, count)
	}
}

// Tsdb定时任务, 将数据通过api发送到tsdb
func forward2TsdbTask(concurrent int) {
	batch := g.Config().Tsdb.Batch // 一次发送,最多batch条数据
	retry := g.Config().Tsdb.MaxRetry
	sema := nsema.NewSemaphore(concurrent)

	for {
		items := TsdbQueue.PopBackBy(batch)
		if len(items) == 0 {
			time.Sleep(DefaultSendTaskSleepInterval)
			continue
		}
		//  同步Call + 有限并发 进行发送
		sema.Acquire()
		go func(itemList []interface{}) {
			defer sema.Release()

			var tsdbBuffer bytes.Buffer
			count := len(itemList)
			for i := 0; i < count; i++ {
				tsdbItem := itemList[i].(*cmodel.TsdbItem)
				tsdbBuffer.WriteString(tsdbItem.TsdbString())
				tsdbBuffer.WriteString("\n")
			}

			var err error
			for i := 0; i < retry; i++ {
				err = TsdbConnPoolHelper.Send(tsdbBuffer.Bytes())
				if err == nil {
					proc.SendToTsdbCnt.IncrBy(int64(len(itemList)))
					break
				}
				time.Sleep(100 * time.Millisecond)
			}

			if err != nil {
				proc.SendToTsdbFailCnt.IncrBy(int64(len(itemList)))
				log.Println(err)
				return
			}
		}(items)
	}
}

func forward2TransferTask(Q *list.SafeListLimited, concurrent int) {
	cfg := g.Config()
	batch := int(cfg.Transfer.Batch)
	maxConns := int64(cfg.Transfer.MaxConns)
	retry := int(cfg.Transfer.Retry)
	if retry < 1 {
		retry = 1
	}

	sema := nsema.NewSemaphore(int(concurrent))
	transNum := len(TransferHostnames)

	for {
		items := Q.PopBackBy(batch)
		count := len(items)
		if count == 0 {
			time.Sleep(time.Millisecond * 50)
			continue
		}

		transItems := make([]*cmodel.MetricValue, count)
		for i := 0; i < count; i++ {
			transItems[i] = convert(items[i].(*cmodel.MetaData))
		}

		sema.Acquire()
		go func(transItems []*cmodel.MetricValue, count int) {
			defer sema.Release()
			var err error

			// 随机遍历transfer列表，直到数据发送成功 或者 遍历完;随机遍历，可以缓解慢transfer
			resp := &cmodel.TransferResponse{}
			sendOk := false

			for j := 0; j < retry && !sendOk; j++ {
				rint := rand.Int()
				for i := 0; i < transNum && !sendOk; i++ {
					idx := (i + rint) % transNum
					host := TransferHostnames[idx]
					addr := TransferMap[host]

					// 过滤掉建连缓慢的host, 否则会严重影响发送速率
					cc := pfc.GetCounterCount(host)
					if cc >= maxConns {
						continue
					}

					pfc.Counter(host, 1)
					err = TransferConnPools.Call(addr, "Transfer.Update", transItems, resp)
					pfc.Counter(host, -1)

					if err == nil {
						sendOk = true
						// statistics
						proc.TransferSendCnt[host].IncrBy(int64(count))
					} else {
						// statistics
						log.Printf("transfer update fail, items size:%d, error:%v, resp:%v", len(transItems), err, resp)
						proc.TransferSendFailCnt[host].IncrBy(int64(count))
					}
				}
			}

			// statistics
			if !sendOk {
				pfc.Meter("SendFail", int64(count))
			} else {
				pfc.Meter("Send", int64(count))
			}
		}(transItems, count)
	}
}

func convert(v *cmodel.MetaData) *cmodel.MetricValue {
	return &cmodel.MetricValue{
		Metric:    v.Metric,
		Endpoint:  v.Endpoint,
		Timestamp: v.Timestamp,
		Step:      v.Step,
		Type:      v.CounterType,
		Tags:      cutils.SortedTags(v.Tags),
		Value:     v.Value,
	}
}
