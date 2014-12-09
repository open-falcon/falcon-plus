package funcs

import (
	"fmt"
	"github.com/open-falcon/agent/g"
	"time"
)

func p(item *g.MetricValue) {
	fmt.Printf("%s=%v\n", item.Metric, item.Value)
}

func PrintAll() {
	for _, item := range AgentMetrics() {
		p(item)
	}

	for _, item := range KernelMetrics() {
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
