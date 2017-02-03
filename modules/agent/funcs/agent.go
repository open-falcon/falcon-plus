package funcs

import (
	"github.com/open-falcon/falcon-plus/common/model"
)

func AgentMetrics() []*model.MetricValue {
	return []*model.MetricValue{GaugeValue("agent.alive", 1)}
}
