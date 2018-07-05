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
	"log"
	"strings"
	"time"

	"github.com/open-falcon/falcon-plus/modules/transfer/g"
	"github.com/open-falcon/falcon-plus/modules/transfer/proc"
	"github.com/toolkits/container/list"
)

const (
	DefaultProcCronPeriod = time.Duration(5) * time.Second    //ProcCron的周期,默认1s
	DefaultLogCronPeriod  = time.Duration(3600) * time.Second //LogCron的周期,默认300s
)

// send_cron程序入口
func startSenderCron() {
	go startProcCron()
	go startLogCron()
}

func startProcCron() {
	for {
		time.Sleep(DefaultProcCronPeriod)
		refreshSendingCacheSize()
	}
}

func startLogCron() {
	for {
		time.Sleep(DefaultLogCronPeriod)
		logConnPoolsProc()
	}
}

func refreshSendingCacheSize() {
	proc.JudgeQueuesCnt.SetCnt(calcSendCacheSize(JudgeQueues))
	proc.GraphQueuesCnt.SetCnt(calcSendCacheSize(GraphQueues))

	cfg := g.Config()
	if cfg.Tsdb.Enabled && TsdbQueue != nil {
		proc.TsdbQueuesCnt.SetCnt(int64(TsdbQueue.Len()))
	}

	if cfg.Influxdb.Enabled && InfluxdbQueue != nil {
		proc.InfluxdbQueuesCnt.SetCnt(int64(InfluxdbQueue.Len()))
	}
}
func calcSendCacheSize(mapList map[string]*list.SafeListLimited) int64 {
	var cnt int64 = 0
	for _, list := range mapList {
		if list != nil {
			cnt += int64(list.Len())
		}
	}
	return cnt
}

func logConnPoolsProc() {
	log.Printf("connPools proc: \n%v", strings.Join(GraphConnPools.Proc(), "\n"))
}
