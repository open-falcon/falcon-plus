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

	"github.com/juju/errors"
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
	influxdbConcurrent := cfg.Influxdb.MaxConns
	p8sRelayConcurrent := cfg.P8sRelay.MaxConns

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

	if influxdbConcurrent < 1 {
		influxdbConcurrent = 1
	}

	if p8sRelayConcurrent < 1 {
		p8sRelayConcurrent = 1
	}

	// init send go-routines
	for node := range cfg.Judge.Cluster {
		queue := JudgeQueues[node]
		go forward2JudgeTask(queue, node, judgeConcurrent)
	}

	for node, nitem := range cfg.Graph.ClusterList {
		for _, addr := range nitem.Addrs {
			queue := GraphQueues[node+addr]
			go forward2GraphTask(queue, node, addr, graphConcurrent)
		}
	}

	if cfg.Tsdb.Enabled {
		go forward2TsdbTask(tsdbConcurrent)
	}

	if cfg.Transfer.Enabled {
		concurrent := transferConcurrent * len(cfg.Transfer.Cluster)
		go forward2TransferTask(TransferQueue, concurrent)
	}

	if cfg.Influxdb.Enabled {
		go forward2InfluxdbTask(influxdbConcurrent)
	}

	if cfg.P8sRelay.Enabled {
		for node, nitem := range cfg.P8sRelay.ClusterList {
			for _, addr := range nitem.Addrs {
				queue := P8sRelayQueues[node+addr]
				go forward2P8sRelayTask(queue, node, addr, p8sRelayConcurrent)
			}
		}
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

func forward2P8sRelayTask(Q *list.SafeListLimited, node string, addr string, concurrent int) {
	batch := g.Config().P8sRelay.Batch
	sema := nsema.NewSemaphore(concurrent)
	for {
		items := Q.PopBackBy(batch)
		count := len(items)
		if count == 0 {
			time.Sleep(DefaultSendTaskSleepInterval)
			continue
		}
		p8sItems := make([]*cmodel.P8sItem, count)
		for i := 0; i < count; i++ {
			p8sItems[i] = items[i].(*cmodel.P8sItem)
		}

		sema.Acquire()
		go func(addr string, p8sItems []*cmodel.P8sItem, count int) {
			defer sema.Release()
			resp := &cmodel.SimpleRpcResponse{}
			var err error
			sendOk := false
			for i := 0; i < 3; i++ {
				err = P8sRelayConnPools.Call(addr, "P8sRelay.Send", p8sItems, resp)
				if err == nil {
					sendOk = true
					break
				}
				time.Sleep(time.Millisecond * 10)
			}
			if !sendOk {
				log.Printf("send p8s relay %s:%s fail: %v", node, addr, err)
				proc.SendToP8sRelayFailCnt.IncrBy(int64(count))
			} else {
				proc.SendToP8sRelayCnt.IncrBy(int64(count))
			}
		}(addr, p8sItems, count)
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

// Transfer定时任务, 将Transfer发送缓存中的数据 通过rpc连接池 发送到Transfer(此时transfer仅仅起到转发数据的功能)
func forward2TransferTask(Q *list.SafeListLimited, concurrent int) {
	cfg := g.Config()
	batch := cfg.Transfer.Batch // 一次发送,最多batch条数据
	maxConns := int64(cfg.Transfer.MaxConns)
	retry := cfg.Transfer.MaxRetry //最多尝试发送retry次
	if retry < 1 {
		retry = 1
	}

	sema := nsema.NewSemaphore(concurrent)
	transNum := len(TransferHostnames)

	for {
		items := Q.PopBackBy(batch)
		count := len(items)
		if count == 0 {
			time.Sleep(DefaultSendTaskSleepInterval)
			continue
		}

		transferItems := make([]*cmodel.MetricValue, count)
		for i := 0; i < count; i++ {
			transferItems[i] = convert(items[i].(*cmodel.MetaData))
		}

		//	同步Call + 有限并发 进行发送
		sema.Acquire()
		go func(transferItems []*cmodel.MetricValue, count int) {
			defer sema.Release()

			// 随机遍历transfer列表，直到数据发送成功 或者 遍历完;随机遍历，可以缓解慢transfer
			resp := &cmodel.TransferResponse{}
			var err error
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
					err = TransferConnPools.Call(addr, "Transfer.Update", transferItems, resp)
					pfc.Counter(host, -1)

					// statistics
					if err == nil {
						sendOk = true
						proc.SendToTransferCnt.IncrBy(int64(count))
					} else {
						log.Printf("transfer update fail, transfer hostname: %s, transfer instance: %s, items size:%d, error:%v, resp:%v", host, addr, len(transferItems), err, resp)
						proc.SendToTransferFailCnt.IncrBy(int64(count))
					}
				}
			}
		}(transferItems, count)
	}
}

// cmodel.MetaData --> cmodel.MetricValue
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

func forward2InfluxdbTask(concurrent int) {
	batch := g.Config().Influxdb.Batch // 一次发送,最多batch条数据
	retry := g.Config().Influxdb.MaxRetry
	sema := nsema.NewSemaphore(concurrent)

	for {
		items := InfluxdbQueue.PopBackBy(batch)
		if len(items) == 0 {
			time.Sleep(DefaultSendTaskSleepInterval)
			continue
		}
		//  同步Call + 有限并发 进行发送
		sema.Acquire()
		go func(itemList []interface{}) {
			defer sema.Release()

			var err error
			c, err := NewInfluxdbClient()
			defer c.Client.Close()

			if err != nil {
				log.Println(errors.ErrorStack(err))
				return
			}

			items := make([]*cmodel.InfluxdbItem, 0, batch)
			for _, i := range itemList {
				items = append(items, i.(*cmodel.InfluxdbItem))
			}

			for i := 0; i < retry; i++ {
				err = c.Send(items)
				if err == nil {
					proc.SendToInfluxdbCnt.IncrBy(int64(len(itemList)))
					break
				}
			}

			if err != nil {
				proc.SendToInfluxdbFailCnt.IncrBy(int64(len(itemList)))
				log.Println(err)
				return
			}
		}(items)
	}
}
