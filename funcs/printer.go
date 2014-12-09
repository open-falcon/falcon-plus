package funcs

import (
	"fmt"
	"github.com/open-falcon/agent/g"
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
}
