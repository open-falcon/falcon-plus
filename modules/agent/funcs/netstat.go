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

package funcs

import (
	"github.com/open-falcon/falcon-plus/common/model"
	"github.com/toolkits/nux"
	"log"
)

var USES = map[string]struct{}{
	"PruneCalled":        {},
	"LockDroppedIcmps":   {},
	"ArpFilter":          {},
	"TW":                 {},
	"DelayedACKLocked":   {},
	"ListenOverflows":    {},
	"ListenDrops":        {},
	"TCPPrequeueDropped": {},
	"TCPTSReorder":       {},
	"TCPDSACKUndo":       {},
	"TCPLoss":            {},
	"TCPLostRetransmit":  {},
	"TCPLossFailures":    {},
	"TCPFastRetrans":     {},
	"TCPTimeouts":        {},
	"TCPSchedulerFailed": {},
	"TCPAbortOnMemory":   {},
	"TCPAbortOnTimeout":  {},
	"TCPAbortFailed":     {},
	"TCPMemoryPressures": {},
	"TCPSpuriousRTOs":    {},
	"TCPBacklogDrop":     {},
	"TCPMinTTLDrop":      {},
}

func NetstatMetrics() (L []*model.MetricValue) {
	tcpExts, err := nux.Netstat("TcpExt")

	if err != nil {
		log.Println(err)
		return
	}

	cnt := len(tcpExts)
	if cnt == 0 {
		return
	}

	for key, val := range tcpExts {
		if _, ok := USES[key]; !ok {
			continue
		}
		L = append(L, CounterValue("TcpExt."+key, val))
	}

	return
}
