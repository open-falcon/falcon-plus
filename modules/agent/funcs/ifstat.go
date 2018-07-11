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
	"github.com/open-falcon/falcon-plus/modules/agent/g"
	"github.com/toolkits/nux"
	"log"
)

func NetMetrics() []*model.MetricValue {
	return CoreNetMetrics(g.Config().Collector.IfacePrefix)
}

func CoreNetMetrics(ifacePrefix []string) []*model.MetricValue {

	netIfs, err := nux.NetIfs(ifacePrefix)
	if err != nil {
		log.Println(err)
		return []*model.MetricValue{}
	}

	cnt := len(netIfs)
	ret := make([]*model.MetricValue, cnt*26)

	for idx, netIf := range netIfs {
		iface := "iface=" + netIf.Iface
		ret[idx*26+0] = CounterValue("net.if.in.bytes", netIf.InBytes, iface)
		ret[idx*26+1] = CounterValue("net.if.in.packets", netIf.InPackages, iface)
		ret[idx*26+2] = CounterValue("net.if.in.errors", netIf.InErrors, iface)
		ret[idx*26+3] = CounterValue("net.if.in.dropped", netIf.InDropped, iface)
		ret[idx*26+4] = CounterValue("net.if.in.fifo.errs", netIf.InFifoErrs, iface)
		ret[idx*26+5] = CounterValue("net.if.in.frame.errs", netIf.InFrameErrs, iface)
		ret[idx*26+6] = CounterValue("net.if.in.compressed", netIf.InCompressed, iface)
		ret[idx*26+7] = CounterValue("net.if.in.multicast", netIf.InMulticast, iface)
		ret[idx*26+8] = CounterValue("net.if.out.bytes", netIf.OutBytes, iface)
		ret[idx*26+9] = CounterValue("net.if.out.packets", netIf.OutPackages, iface)
		ret[idx*26+10] = CounterValue("net.if.out.errors", netIf.OutErrors, iface)
		ret[idx*26+11] = CounterValue("net.if.out.dropped", netIf.OutDropped, iface)
		ret[idx*26+12] = CounterValue("net.if.out.fifo.errs", netIf.OutFifoErrs, iface)
		ret[idx*26+13] = CounterValue("net.if.out.collisions", netIf.OutCollisions, iface)
		ret[idx*26+14] = CounterValue("net.if.out.carrier.errs", netIf.OutCarrierErrs, iface)
		ret[idx*26+15] = CounterValue("net.if.out.compressed", netIf.OutCompressed, iface)
		ret[idx*26+16] = CounterValue("net.if.total.bytes", netIf.TotalBytes, iface)
		ret[idx*26+17] = CounterValue("net.if.total.packets", netIf.TotalPackages, iface)
		ret[idx*26+18] = CounterValue("net.if.total.errors", netIf.TotalErrors, iface)
		ret[idx*26+19] = CounterValue("net.if.total.dropped", netIf.TotalDropped, iface)
		ret[idx*26+20] = GaugeValue("net.if.speed.bits", netIf.SpeedBits, iface)
		ret[idx*26+21] = CounterValue("net.if.in.percent", netIf.InPercent, iface)
		ret[idx*26+22] = CounterValue("net.if.out.percent", netIf.OutPercent, iface)
		ret[idx*26+23] = CounterValue("net.if.in.bits", netIf.InBytes*8, iface)
		ret[idx*26+24] = CounterValue("net.if.out.bits", netIf.OutBytes*8, iface)
		ret[idx*26+25] = CounterValue("net.if.total.bits", netIf.TotalBytes*8, iface)
	}
	return ret
}
