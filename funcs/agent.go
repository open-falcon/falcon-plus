package funcs

import (
	"github.com/open-falcon/agent/g"
)

func AgentMetrics() []*g.MetricValue {
	return []*g.MetricValue{GaugeValue("agent.alive", 1)}
}
