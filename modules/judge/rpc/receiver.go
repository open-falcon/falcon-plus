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
	"log"
	"time"

	cmodel "github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/judge/g"
	"github.com/open-falcon/falcon-plus/modules/judge/model"
	"github.com/open-falcon/falcon-plus/modules/judge/store"
	"github.com/open-falcon/falcon-plus/modules/judge/string_matcher"
)

type Judge int

func (this *Judge) Ping(req cmodel.NullRpcRequest, resp *cmodel.SimpleRpcResponse) error {
	return nil
}

func (this *Judge) Send(items []*cmodel.JudgeItem, resp *cmodel.SimpleRpcResponse) error {
	cfg := g.Config()
	remain := cfg.Remain
	// 把当前时间的计算放在最外层，是为了减少获取时间时的系统调用开销
	now := time.Now().Unix()

	for _, item := range items {
		exists := g.FilterMap.Exists(item.Metric)
		if !exists {
			continue
		}

		if item.Metric != "str.match" {
			pk := item.PrimaryKey()
			store.HistoryBigMap[pk[0:2]].PushFrontAndMaintain(pk, item, remain, now)

		} else if item.Metric == "str.match" && item.JudgeType == g.GAUGE && item.ValueRaw != "" {
			yesEndpoint := g.StrMatcherMap.Match(item.Endpoint, item.ValueRaw)
			yesTag := g.StrMatcherExpMap.Match(item.Tags, item.ValueRaw)

			if yesEndpoint || yesTag {
				pk := item.PrimaryKey()
				store.HistoryBigMap[pk[0:2]].PushFrontAndMaintain(pk, item, remain, now)

				// save matched string into SQL DB
				if cfg.StringMatcher.Enabled {
					success := string_matcher.Producer.Append(item)
					if !success {
						log.Println("string_matcher.Producer failed")
					}
				}
			}
		}

	}
	return nil
}

func (this *Judge) SendE(items []*cmodel.EMetric, resp *cmodel.SimpleRpcResponse) error {
	cfg := g.Config()

	// remain default is 11, why?
	remain := cfg.Remain
	// 把当前时间的计算放在最外层，是为了减少获取时间时的系统调用开销
	now := time.Now().Unix()

	for _, item := range items {
		exists := g.EFilterMap.Exists(item.Metric)
		if !exists {
			continue
		}

		pksum := utils.Md5(item.PK())
		model.EHistoryBigMap[pksum[0:2]].PushFrontAndMaintain(pksum, item, remain, now)

	}
	return nil
}
