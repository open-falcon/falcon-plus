package funcs

import (
	"github.com/open-falcon/agent/g"
	"github.com/toolkits/nux"
	"log"
)

func NetMetrics() []*g.MetricValue {
	return CoreNetMetrics(g.Config().Collector.IfacePrefix)
}

func CoreNetMetrics(ifacePrefix []string) []*g.MetricValue {

	netIfs, err := nux.NetIfs(ifacePrefix)
	if err != nil {
		log.Println(err)
		return []*g.MetricValue{}
	}

	cnt := len(netIfs)
	ret := make([]*g.MetricValue, cnt*20)

	for idx, netIf := range netIfs {
		iface := "iface=" + netIf.Iface
		ret[idx*20+0] = CounterValue("net.if.in.bytes", netIf.InBytes, iface)
		ret[idx*20+1] = CounterValue("net.if.in.packets", netIf.InPackages, iface)
		ret[idx*20+2] = CounterValue("net.if.in.errors", netIf.InErrors, iface)
		ret[idx*20+3] = CounterValue("net.if.in.dropped", netIf.InDropped, iface)
		ret[idx*20+4] = CounterValue("net.if.in.fifo.errs", netIf.InFifoErrs, iface)
		ret[idx*20+5] = CounterValue("net.if.in.frame.errs", netIf.InFrameErrs, iface)
		ret[idx*20+6] = CounterValue("net.if.in.compressed", netIf.InCompressed, iface)
		ret[idx*20+7] = CounterValue("net.if.in.multicast", netIf.InMulticast, iface)
		ret[idx*20+8] = CounterValue("net.if.out.bytes", netIf.OutBytes, iface)
		ret[idx*20+9] = CounterValue("net.if.out.packets", netIf.OutPackages, iface)
		ret[idx*20+10] = CounterValue("net.if.out.errors", netIf.OutErrors, iface)
		ret[idx*20+11] = CounterValue("net.if.out.dropped", netIf.OutDropped, iface)
		ret[idx*20+12] = CounterValue("net.if.out.fifo.errs", netIf.OutFifoErrs, iface)
		ret[idx*20+13] = CounterValue("net.if.out.collisions", netIf.OutCollisions, iface)
		ret[idx*20+14] = CounterValue("net.if.out.carrier.errs", netIf.OutCarrierErrs, iface)
		ret[idx*20+15] = CounterValue("net.if.out.compressed", netIf.OutCompressed, iface)
		ret[idx*20+16] = CounterValue("net.if.total.bytes", netIf.TotalBytes, iface)
		ret[idx*20+17] = CounterValue("net.if.total.packets", netIf.TotalPackages, iface)
		ret[idx*20+18] = CounterValue("net.if.total.errors", netIf.TotalErrors, iface)
		ret[idx*20+19] = CounterValue("net.if.total.dropped", netIf.TotalDropped, iface)
	}
	return ret
}
