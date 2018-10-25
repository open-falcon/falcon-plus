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
	log "github.com/Sirupsen/logrus"
)

func MemMetrics() []*model.MetricValue {
	m, err := nux.MemInfo()
	if err != nil {
		log.Println(err)
		return nil
	}

	memFree := m.MemFree + m.Buffers + m.Cached
	memUsed := m.MemTotal - memFree

	pmemFree := 0.0
	pmemUsed := 0.0
	if m.MemTotal != 0 {
		pmemFree = float64(memFree) * 100.0 / float64(m.MemTotal)
		pmemUsed = float64(memUsed) * 100.0 / float64(m.MemTotal)
	}

	pswapFree := 0.0
	pswapUsed := 0.0
	if m.SwapTotal != 0 {
		pswapFree = float64(m.SwapFree) * 100.0 / float64(m.SwapTotal)
		pswapUsed = float64(m.SwapUsed) * 100.0 / float64(m.SwapTotal)
	}

	return []*model.MetricValue{
		GaugeValue("mem.memtotal", m.MemTotal),
		GaugeValue("mem.memused", memUsed),
		GaugeValue("mem.memfree", memFree),
		GaugeValue("mem.swaptotal", m.SwapTotal),
		GaugeValue("mem.swapused", m.SwapUsed),
		GaugeValue("mem.swapfree", m.SwapFree),
		GaugeValue("mem.memfree.percent", pmemFree),
		GaugeValue("mem.memused.percent", pmemUsed),
		GaugeValue("mem.swapfree.percent", pswapFree),
		GaugeValue("mem.swapused.percent", pswapUsed),
	}

}
