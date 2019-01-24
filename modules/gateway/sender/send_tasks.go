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
	log "github.com/Sirupsen/logrus"
	"github.com/open-falcon/falcon-plus/common/encrypt"
	"math/rand"
	"time"

	pfc "github.com/niean/goperfcounter"
	cmodel "github.com/open-falcon/falcon-plus/common/model"
	cutils "github.com/open-falcon/falcon-plus/common/utils"
	nsema "github.com/toolkits/concurrent/semaphore"
	nlist "github.com/toolkits/container/list"

	"github.com/open-falcon/falcon-plus/modules/gateway/g"
)

func startSendTasks() {
	cfg := g.Config()
	concurrent := cfg.Transfer.MaxConns * int32(len(cfg.Transfer.Cluster))
	go forward2TransferTask(SenderQueue, concurrent)
}

func forward2TransferTask(Q *nlist.SafeListLimited, concurrent int32) {
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
			aesKey := g.Config().AesKey

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

					if len(aesKey) > 16 {
						encrypted, _ := encrypt.Encrypt(encrypt.Encode(transItems), []byte(g.Config().AesKey))
						err = SenderConnPools.Call(addr, "Transfer.Update", []cmodel.Encrypt{{Byte: encrypted}}, resp)
					} else {
						err = SenderConnPools.Call(addr, "Transfer.Update", transItems, resp)
					}

					pfc.Counter(host, -1)

					if err == nil {
						sendOk = true
						// statistics
						TransferSendCnt[host].IncrBy(int64(count))
					} else {
						// statistics
						log.Errorf("transfer update fail, items size:%d, error:%v, resp:%v", len(transItems), err, resp)
						TransferSendFailCnt[host].IncrBy(int64(count))
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
