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

package rpc

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	cmodel "github.com/open-falcon/falcon-plus/common/model"
	cutils "github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/transfer/g"
	"github.com/open-falcon/falcon-plus/modules/transfer/proc"
	"github.com/open-falcon/falcon-plus/modules/transfer/sender"
)

var (
	NOT_FOUND = -1
)

type Transfer int

type TransferResp struct {
	Msg        string
	Total      int
	ErrInvalid int
	Latency    int64
}

func (t *TransferResp) String() string {
	s := fmt.Sprintf("TransferResp total=%d, err_invalid=%d, latency=%dms",
		t.Total, t.ErrInvalid, t.Latency)
	if t.Msg != "" {
		s = fmt.Sprintf("%s, msg=%s", s, t.Msg)
	}
	return s
}

func (this *Transfer) Ping(req cmodel.NullRpcRequest, resp *cmodel.SimpleRpcResponse) error {
	return nil
}

func (t *Transfer) Update(args []*cmodel.MetricValue, reply *cmodel.TransferResponse) error {
	return RecvMetricValues(args, reply, "rpc")
}

// process new metric values
func RecvMetricValues(args []*cmodel.MetricValue, reply *cmodel.TransferResponse, from string) error {
	cfg := g.Config()

	start := time.Now()
	reply.Invalid = 0

	errmsg := []string{}

	items := []*cmodel.MetaData{}

	for _, v := range args {
		if cfg.Debug {
			log.Println("metric", v)
		}

		if v == nil {
			errmsg = append(errmsg, "metric empty")
			reply.Invalid += 1
			continue
		}

		// 历史遗留问题.
		// 老版本agent上报的metric=kernel.hostname的数据,其取值为string类型,现在已经不支持了;所以,这里硬编码过滤掉
		if v.Metric == "kernel.hostname" {
			errmsg = append(errmsg, "skip metric=kernel.hostname")
			reply.Invalid += 1
			continue
		}

		if v.Metric == "" || v.Endpoint == "" {
			errmsg = append(errmsg, "metric or endpoint is empty")
			reply.Invalid += 1
			continue
		}

		if v.Type != g.COUNTER && v.Type != g.GAUGE && v.Type != g.DERIVE && v.Type != g.STRMATCH {
			reply.Invalid += 1
			errmsg = append(errmsg, "got unexpected counterType"+v.Type)
			continue
		}

		if v.Value == "" {
			errmsg = append(errmsg, v.Metric+" value is empty")
			reply.Invalid += 1
			continue
		}

		if v.Step <= 0 {
			errmsg = append(errmsg, "Step <= 0")
			reply.Invalid += 1
			continue
		}

		if len(v.Metric)+len(v.Tags) > 510 {
			errmsg = append(errmsg, " Metric+Tags too long")
			reply.Invalid += 1
			continue
		}

		now := start.Unix()
		if v.Timestamp <= 0 || v.Timestamp > now*2 {
			errmsg = append(errmsg, v.Metric+" Timestamp invalid")
			v.Timestamp = now
		}

		fv := &cmodel.MetaData{
			Metric:      v.Metric,
			Endpoint:    v.Endpoint,
			Timestamp:   v.Timestamp,
			Step:        v.Step,
			CounterType: v.Type,
			Tags:        cutils.DictedTagstring(v.Tags), //TODO tags键值对的个数,要做一下限制
		}

		valid := true
		var vv float64
		var err error

		if v.Type != g.STRMATCH {
			switch cv := v.Value.(type) {
			case string:
				vv, err = strconv.ParseFloat(cv, 64)
				if err != nil {
					valid = false
				}

			case float64:
				vv = cv
			case int64:
				vv = float64(cv)
			default:
				valid = false
			}
		} else {
			switch v.Value.(type) {
			case string:
				fv.ValueRaw = v.Value.(string)
				vv = float64(1.0)
			default:
				valid = false
			}
		}

		if !valid {
			errmsg = append(errmsg, "parse value into float64/string failed")
			reply.Invalid += 1
			continue
		}

		fv.Value = vv

		items = append(items, fv)
	}

	// statistics
	cnt := int64(len(items))
	proc.RecvCnt.IncrBy(cnt)
	if from == "rpc" {
		proc.RpcRecvCnt.IncrBy(cnt)
	} else if from == "http" {
		proc.HttpRecvCnt.IncrBy(cnt)
	}

	if cfg.Graph.Enabled {
		sender.Push2GraphSendQueue(items)
	}

	if cfg.Judge.Enabled {
		sender.Push2JudgeSendQueue(items)
	}

	if cfg.Tsdb.Enabled {
		sender.Push2TsdbSendQueue(items)
	}
	if reply.Invalid == 0 {
		reply.Message = "ok"
	} else {
		reply.Message = strings.Join(errmsg, ";\n")
	}
	reply.Total = len(args)
	reply.Latency = (time.Now().UnixNano() - start.UnixNano()) / 1000000

	return nil
}
