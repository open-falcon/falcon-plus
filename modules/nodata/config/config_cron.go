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

package config

import (
	"log"
	"time"

	"github.com/toolkits/container/nmap"
	tcron "github.com/toolkits/cron"
	ttime "github.com/toolkits/time"

	cutils "github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/nodata/config/service"
	"github.com/open-falcon/falcon-plus/modules/nodata/g"
)

var (
	ndconfigCron     = tcron.New()
	ndconfigCronSpec = "50 */2 * * * ?"
)

func StartNdConfigCron() {
	ndconfigCron.AddFuncCC(ndconfigCronSpec, func() {
		start := time.Now().Unix()
		cnt, _ := syncNdConfig()
		end := time.Now().Unix()
		if g.Config().Debug {
			log.Printf("config cron, cnt %d, time %ds, start %s\n", cnt, end-start, ttime.FormatTs(start))
		}

		// statistics
		g.ConfigCronCnt.Incr()
		g.ConfigLastTs.SetCnt(end - start)
		g.ConfigLastCnt.SetCnt(int64(cnt))
	}, 1)
	ndconfigCron.Start()
}

func SyncNdConfigOnce() int {
	cnt, _ := syncNdConfig()
	return cnt
}

func syncNdConfig() (cnt int, errt error) {
	// get configs
	configs := service.GetMockCfgFromDB()
	// restruct
	nm := nmap.NewSafeMap()
	for _, ndc := range configs {
		endpoint := ndc.Endpoint
		metric := ndc.Metric
		tags := ndc.Tags
		if endpoint == "" {
			log.Printf("bad config: %+v\n", ndc)
			continue
		}
		pk := cutils.PK(endpoint, metric, tags)
		nm.Put(pk, ndc)
	}

	// cache
	SetNdConfigMap(nm)

	return nm.Size(), nil
}
