package funcs

import (
	"github.com/open-falcon/falcon-plus/modules/agent/g"
	"github.com/open-falcon/falcon-plus/common/model"
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
	ret := make([]*model.MetricValue, cnt*23)

	for idx, netIf := range netIfs {
		iface := "iface=" + netIf.Iface
		ret[idx*23+0] = CounterValue("net.if.in.bytes", netIf.InBytes, iface)
		ret[idx*23+1] = CounterValue("net.if.in.packets", netIf.InPackages, iface)
		ret[idx*23+2] = CounterValue("net.if.in.errors", netIf.InErrors, iface)
		ret[idx*23+3] = CounterValue("net.if.in.dropped", netIf.InDropped, iface)
		ret[idx*23+4] = CounterValue("net.if.in.fifo.errs", netIf.InFifoErrs, iface)
		ret[idx*23+5] = CounterValue("net.if.in.frame.errs", netIf.InFrameErrs, iface)
		ret[idx*23+6] = CounterValue("net.if.in.compressed", netIf.InCompressed, iface)
		ret[idx*23+7] = CounterValue("net.if.in.multicast", netIf.InMulticast, iface)
		ret[idx*23+8] = CounterValue("net.if.out.bytes", netIf.OutBytes, iface)
		ret[idx*23+9] = CounterValue("net.if.out.packets", netIf.OutPackages, iface)
		ret[idx*23+10] = CounterValue("net.if.out.errors", netIf.OutErrors, iface)
		ret[idx*23+11] = CounterValue("net.if.out.dropped", netIf.OutDropped, iface)
		ret[idx*23+12] = CounterValue("net.if.out.fifo.errs", netIf.OutFifoErrs, iface)
		ret[idx*23+13] = CounterValue("net.if.out.collisions", netIf.OutCollisions, iface)
		ret[idx*23+14] = CounterValue("net.if.out.carrier.errs", netIf.OutCarrierErrs, iface)
		ret[idx*23+15] = CounterValue("net.if.out.compressed", netIf.OutCompressed, iface)
		ret[idx*23+16] = CounterValue("net.if.total.bytes", netIf.TotalBytes, iface)
		ret[idx*23+17] = CounterValue("net.if.total.packets", netIf.TotalPackages, iface)
		ret[idx*23+18] = CounterValue("net.if.total.errors", netIf.TotalErrors, iface)
		ret[idx*23+19] = CounterValue("net.if.total.dropped", netIf.TotalDropped, iface)
		ret[idx*23+20] = GaugeValue("net.if.speed.bits", netIf.SpeedBits, iface)
		ret[idx*23+21] = CounterValue("net.if.in.percent", netIf.InPercent, iface)
		ret[idx*23+22] = CounterValue("net.if.out.percent", netIf.OutPercent, iface)
	}
	return ret
}
