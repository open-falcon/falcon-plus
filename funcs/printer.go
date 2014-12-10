package funcs

import (
	"fmt"
	"github.com/open-falcon/agent/g"
	"time"
)

func p(item *g.MetricValue) {
	fmt.Printf("%s=%v[tags:%s]\n", item.Metric, item.Value, item.Tags)
}

func PrintAll() {
	for _, item := range AgentMetrics() {
		p(item)
	}

	for _, item := range KernelMetrics() {
		p(item)
	}

	for _, item := range DeviceMetrics() {
		p(item)
	}

	err := UpdateCpuStat()
	if err != nil {
		fmt.Println(err)
	}

	time.Sleep(time.Second)
	UpdateCpuStat()

	for _, item := range CpuMetrics() {
		p(item)
	}
}
