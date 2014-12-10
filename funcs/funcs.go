package funcs

import (
	"github.com/open-falcon/agent/g"
)

type FuncsAndInterval struct {
	Fs       []func() []*g.MetricValue
	Interval int
}

var Mappers []FuncsAndInterval

func BuildMappers() {
	interval := g.Config().Transfer.Interval
	Mappers = []FuncsAndInterval{
		FuncsAndInterval{
			Fs: []func() []*g.MetricValue{
				AgentMetrics,
				CpuMetrics,
				NetMetrics,
				KernelMetrics,
				LoadAvgMetrics,
				MemMetrics,
				DiskIOMetrics,
				IOStatsMetrics,
				NetstatMetrics,
				ProcMetrics,
			},
			Interval: interval,
		},
		FuncsAndInterval{
			Fs: []func() []*g.MetricValue{
				DeviceMetrics,
			},
			Interval: interval,
		},
		FuncsAndInterval{
			Fs: []func() []*g.MetricValue{
				PortMetrics,
				SocketStatSummaryMetrics,
			},
			Interval: interval,
		},
	}
}
