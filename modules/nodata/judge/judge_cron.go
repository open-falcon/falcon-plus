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

package judge

import (
	"log"
	"time"

	cmodel "github.com/open-falcon/falcon-plus/common/model"
	cutils "github.com/open-falcon/falcon-plus/common/utils"
	tcron "github.com/toolkits/cron"
	ttime "github.com/toolkits/time"

	"github.com/open-falcon/falcon-plus/modules/nodata/collector"
	"github.com/open-falcon/falcon-plus/modules/nodata/config"
	"github.com/open-falcon/falcon-plus/modules/nodata/g"
	"github.com/open-falcon/falcon-plus/modules/nodata/sender"
)

var (
	judgeCron = tcron.New()
	judgeCronSpec = "0 * * * * ?"
)

func StartJudgeCron() {
	judgeCron.AddFuncCC(judgeCronSpec, func() {
		start := time.Now().Unix()
		judge()
		end := time.Now().Unix()
		if g.Config().Debug {
			log.Printf("judge cron, time %ds, start %s\n", end - start, ttime.FormatTs(start))
		}

		// statistics
		g.JudgeCronCnt.Incr()
		g.JudgeLastTs.SetCnt(end - start)

		// trigger sender
		sender.SendMockOnceAsync()
	}, 1)
	judgeCron.Start()
}

// Do Judge
func judge() {
	now := time.Now().Unix()
	keys := config.Keys()
	delayPeriod := int(g.Config().Sender.DelayPeriod)
	for _, key := range keys {
		ndcfg, found := config.GetNdConfig(key)
		if !found {
			//策略不存在,不处理
			continue
		}
		step := ndcfg.Step
		mock := ndcfg.Mock

		//item, found := collector.GetFirstItem(key)
		item, found := collector.GetItemByIndex(key, delayPeriod)
		if !found {
			//没有数据,未开始采集,不处理
			continue
		}

		lastTs := now - getTimeout(step)
		if item.FStatus != "OK" || item.FTs < lastTs {
			//数据采集失败,不处理
			continue
		}

		if fCompare(mock, item.Value) == 0 {
			//采集到的数据为mock数据,则认为上报超时了
			if LastTs(key) + step <= now {
				TurnNodata(key, now)
				genMock(genTs(now, step), key, ndcfg)
			}
			continue
		}

		if item.Ts < lastTs {
			//数据过期, 则认为上报超时
			if LastTs(key) + step <= now {
				TurnNodata(key, now)
				genMock(genTs(now, step), key, ndcfg)
			}
			continue
		}

		TurnOk(key, now)
	}
}

func genMock(ts int64, key string, ndcfg *cmodel.NodataConfig) {
	sender.AddMock(key, ndcfg.Endpoint, ndcfg.Metric, cutils.SortedTags(ndcfg.Tags), ts, ndcfg.Type, ndcfg.Step, ndcfg.Mock)
}

//mock的数据,要前移1+个周期、防止覆盖正常值
func genTs(nowTs int64, step int64) int64 {
	if step < 1 {
		step = 60
	}
	delayPeriod := int64(g.Config().Sender.DelayPeriod)
	return nowTs - nowTs % step - delayPeriod * step
}

func getTimeout(step int64) int64 {
	delayPeriod := int64(g.Config().Sender.DelayPeriod)
	if step < 60 {
		return (delayPeriod + 1) * 60 //60*3
	}

	return step * (delayPeriod + 1)
}

const minfloat64 = 0.000001

func fCompare(left, right float64) int {
	sub := left - right
	if sub > -minfloat64 && sub < minfloat64 {
		return 0
	}
	if sub >= minfloat64 {
		return 1
	}
	return -1
}
