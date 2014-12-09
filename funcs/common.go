package funcs

import (
	"github.com/open-falcon/agent/g"
	"strings"
)

func NewMetricValue(metric string, val interface{}, dataType string, tags ...string) *g.MetricValue {
	mv := g.MetricValue{
		Metric: metric,
		Value:  val,
		Type:   dataType,
	}

	size := len(tags)

	if size > 0 {
		mv.Tags = strings.Join(tags, ",")
	}

	return &mv
}

func GaugeValue(metric string, val interface{}, tags ...string) *g.MetricValue {
	return NewMetricValue(metric, val, "GAUGE", tags...)
}

func CounterValue(metric string, val interface{}, tags ...string) *g.MetricValue {
	return NewMetricValue(metric, val, "COUNTER", tags...)
}
