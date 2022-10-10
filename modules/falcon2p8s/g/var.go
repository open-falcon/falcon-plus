package g

import (
	"sync"
	"time"

	"github.com/oleiade/lane"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	IsScraping               bool
	P8sItemQueue             = lane.NewQueue()
	CollectorMap             sync.Map
	MetricTypeMap            sync.Map
	CounterCollectorValueMap sync.Map
	LastUpdateTimeOfCounter  sync.Map
	LastUpdateTimeOfGauge    sync.Map
)

type LastUpdateTimeItem struct {
	LastUpdateTime time.Time
	PK             string
	TagValues      []string
}

type UncheckedCollector struct {
	C          prometheus.Collector
	MetricType string
}

func (u UncheckedCollector) Describe(_ chan<- *prometheus.Desc) {}
func (u UncheckedCollector) Collect(c chan<- prometheus.Metric) {
	u.C.Collect(c)
}
