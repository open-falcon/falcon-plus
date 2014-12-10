package funcs

import (
	"github.com/open-falcon/agent/g"
	"github.com/toolkits/nux"
	"log"
)

func LoadAvgMetrics() []*g.MetricValue {
	load, err := nux.LoadAvg()
	if err != nil {
		log.Println(err)
		return nil
	}

	return []*g.MetricValue{
		GaugeValue("load.1min", load.Avg1min),
		GaugeValue("load.5min", load.Avg5min),
		GaugeValue("load.15min", load.Avg15min),
	}

}
