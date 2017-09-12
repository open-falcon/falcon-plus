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
	"time"

	pfc "github.com/niean/goperfcounter"
)

const (
	DefaultProcCronPeriod = time.Duration(5) * time.Second //ProcCron的周期,默认1s
)

// send_cron程序入口
func startSenderCron() {
	go startProcCron()
}

func startProcCron() {
	for {
		time.Sleep(DefaultProcCronPeriod)
		refreshSendingCacheSize()
	}
}

func refreshSendingCacheSize() {
	pfc.Gauge("SendQueueSize", int64(SenderQueue.Len()))
}
