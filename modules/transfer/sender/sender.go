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
	"fmt"
	"log"
	"strconv"
	"strings"

	"time"

	"github.com/influxdata/influxdb/client/v2"
	backend "github.com/open-falcon/falcon-plus/common/backend_pool"
	cmodel "github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/transfer/g"
	"github.com/open-falcon/falcon-plus/modules/transfer/proc"
	rings "github.com/toolkits/consistent/rings"
	nlist "github.com/toolkits/container/list"
)

const (
	DefaultSendQueueMaxSize = 102400 //10.24w
)

// 默认参数
var (
	MinStep int //最小上报周期,单位sec
)

// 服务节点的一致性哈希环
// pk -> node
var (
	JudgeNodeRing    *rings.ConsistentHashNodeRing
	GraphNodeRing    *rings.ConsistentHashNodeRing
	P8sRelayNodeRing *rings.ConsistentHashNodeRing
)

// 发送缓存队列
// node -> queue_of_data
var (
	TsdbQueue      *nlist.SafeListLimited
	JudgeQueues    = make(map[string]*nlist.SafeListLimited)
	GraphQueues    = make(map[string]*nlist.SafeListLimited)
	P8sRelayQueues = make(map[string]*nlist.SafeListLimited)
	TransferQueue  *nlist.SafeListLimited
	InfluxdbQueue  *nlist.SafeListLimited
)

// transfer的主机列表，以及主机名和地址的映射关系
// 用于随机遍历transfer
var (
	TransferMap       = make(map[string]string, 0)
	TransferHostnames = make([]string, 0)
)

// 连接池
// node_address -> connection_pool
var (
	JudgeConnPools     *backend.SafeRpcConnPools
	TsdbConnPoolHelper *backend.TsdbConnPoolHelper
	GraphConnPools     *backend.SafeRpcConnPools
	TransferConnPools  *backend.SafeRpcConnPools
	P8sRelayConnPools  *backend.SafeRpcConnPools
)

// infludbConn
type InfluxClient struct {
	Client    client.Client
	Database  string
	Precision string
}

func NewInfluxdbClient() (*InfluxClient, error) {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     g.Config().Influxdb.Address,
		Username: g.Config().Influxdb.Username,
		Password: g.Config().Influxdb.Password,
		Timeout:  time.Millisecond * time.Duration(g.Config().Influxdb.Timeout),
	})

	if err != nil {
		return nil, err
	}

	return &InfluxClient{
		Client:    c,
		Database:  g.Config().Influxdb.Database,
		Precision: g.Config().Influxdb.Precision,
	}, nil
}

func (c *InfluxClient) Send(items []*cmodel.InfluxdbItem) error {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  c.Database,
		Precision: c.Precision,
	})
	if err != nil {
		log.Println("create batch points error: ", err)
		return err
	}

	for _, item := range items {
		pt, err := client.NewPoint(item.Measurement, item.Tags, item.Fileds, time.Unix(item.Timestamp, 0))
		if err != nil {
			log.Println("create new points error: ", err)
			continue
		}
		bp.AddPoint(pt)
	}

	return c.Client.Write(bp)
}

// 初始化数据发送服务, 在main函数中调用
func Start() {
	// 初始化默认参数
	MinStep = g.Config().MinStep
	if MinStep < 1 {
		MinStep = 30 //默认30s
	}
	//
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
		if step < MinStep {
			step = MinStep
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

func Push2P8sRelaySendQueue(items []*cmodel.MetaData) {
	cfg := g.Config().P8sRelay
	notSyncMetrics := g.Config().P8sRelay.NotSyncMetrics
	for _, item := range items {
		// 过滤同步到Prometheus的监控指标
		sync := true
		for _, m := range notSyncMetrics {
			if strings.HasPrefix(item.Metric, m) {
				sync = false
				break
			}
		}
		if !sync {
			continue
		}
		p8sItem, err := convert2P8sRelayItem(item)
		if err != nil {
			log.Println("E:", err)
			continue
		}
		pk := item.PK()
		// statistics
		proc.RecvDataTrace.Trace(pk, item)
		proc.RecvDataFilter.Filter(pk, item.Value, item)

		node, err := P8sRelayNodeRing.GetNode(pk)
		if err != nil {
			log.Println("E:", err)
			continue
		}
		cnode := cfg.ClusterList[node]
		errCnt := 0
		for _, addr := range cnode.Addrs {
			Q := P8sRelayQueues[node+addr]
			if !Q.PushFront(p8sItem) {
				errCnt += 1
			}
		}
		// statistics
		if errCnt > 0 {
			proc.SendToP8sRelayDropCnt.Incr()
		}

	}
}

// 将数据 打入 某个Graph的发送缓存队列, 具体是哪一个Graph 由一致性哈希 决定
func Push2GraphSendQueue(items []*cmodel.MetaData) {
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

		cnode := cfg.ClusterList[node]
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
	}
}

func convert2P8sRelayItem(d *cmodel.MetaData) (*cmodel.P8sItem, error) {
	item := &cmodel.P8sItem{}

	item.Endpoint = d.Endpoint
	item.Metric = d.Metric
	item.Tags = d.Tags
	item.Timestamp = d.Timestamp
	item.Value, _ = strconv.ParseFloat(fmt.Sprintf("%.3f", d.Value), 64)
	item.Step = int(d.Step)
	if item.Step < MinStep {
		item.Step = MinStep
	}

	if d.CounterType == g.GAUGE {
		item.MetricType = g.GAUGE
	} else if d.CounterType == g.COUNTER {
		item.MetricType = g.COUNTER
	} else {
		log.Printf("Error not supported counter type: %s", d.CounterType)
		return item, fmt.Errorf("not_supported_counter_type")
	}

	item.Timestamp = alignTs(item.Timestamp, int64(item.Step)) //item.Timestamp - item.Timestamp%int64(item.Step)

	return item, nil
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
	if item.Step < MinStep {
		item.Step = MinStep
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

func Push2TransferSendQueue(items []*cmodel.MetaData) {
	for _, item := range items {
		isSuccess := TransferQueue.PushFront(item)

		if !isSuccess {
			proc.SendToTransferDropCnt.Incr()
		}
	}
}

// 将原始数据插入到influxdb缓存队列
func Push2InfluxdbSendQueue(items []*cmodel.MetaData) {
	for _, item := range items {
		influxdbItem := convert2InfluxdbItem(item)
		isSuccess := InfluxdbQueue.PushFront(influxdbItem)

		if !isSuccess {
			proc.SendToInfluxdbDropCnt.Incr()
		}
	}
}

func convert2InfluxdbItem(d *cmodel.MetaData) *cmodel.InfluxdbItem {
	t := cmodel.InfluxdbItem{Tags: make(map[string]string), Fileds: make(map[string]interface{})}

	for k, v := range d.Tags {
		t.Tags[k] = v
	}
	t.Tags["endpoint"] = d.Endpoint
	t.Measurement = d.Metric
	t.Fileds["value"] = d.Value
	t.Timestamp = d.Timestamp

	return &t
}
