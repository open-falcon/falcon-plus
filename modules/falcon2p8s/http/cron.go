package http

import (
	"time"

	"github.com/open-falcon/falcon-plus/modules/falcon2p8s/g"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var queueSize = promauto.NewGauge(prometheus.GaugeOpts{
	Name: "falcon2p8s_relay_queue_size",
	Help: "The size of falcon to prometheus relay server",
})

func getqueueSize() {
	for {
		time.Sleep(time.Second * 30)
		queueSize.Set(float64(g.P8sItemQueue.Size()))
	}
}

func cleanOutdatedGaugeMetrics() {
	for {
		time.Sleep(1 * time.Minute)
		g.LastUpdateTimeOfGauge.Range(func(key, value interface{}) bool {
			lastUpdateTimeItem := value.(g.LastUpdateTimeItem)
			if time.Since(lastUpdateTimeItem.LastUpdateTime) > 70*time.Second {
				if counter, ok := g.CollectorMap.Load(lastUpdateTimeItem.PK); ok {
					if uncheckedCollector, ok := counter.(g.UncheckedCollector); ok {
						uncheckedCollector.C.(*prometheus.GaugeVec).DeleteLabelValues(lastUpdateTimeItem.TagValues...)
						g.LastUpdateTimeOfGauge.Delete(key)
					}
				}
			}
			return true
		})
	}
}

func cleanOutdatedCounterMetrics() {
	for {
		time.Sleep(1 * time.Minute)
		g.LastUpdateTimeOfCounter.Range(func(key, value interface{}) bool {
			lastUpdateTimeItem := value.(g.LastUpdateTimeItem)
			if time.Since(lastUpdateTimeItem.LastUpdateTime) > 70*time.Second {
				if counter, ok := g.CollectorMap.Load(lastUpdateTimeItem.PK); ok {
					if uncheckedCollector, ok := counter.(g.UncheckedCollector); ok {
						uncheckedCollector.C.(*prometheus.CounterVec).DeleteLabelValues(lastUpdateTimeItem.TagValues...)
						g.LastUpdateTimeOfCounter.Delete(key)
						g.CounterCollectorValueMap.Delete(key)
					}
				}
			}
			return true
		})
	}
}
