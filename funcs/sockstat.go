package funcs

import (
	"github.com/open-falcon/agent/g"
	"github.com/toolkits/nux"
	"log"
)

func SocketStatSummaryMetrics() (L []*g.MetricValue) {
	ssMap, err := nux.SocketStatSummary()
	if err != nil {
		log.Println(err)
		return
	}

	for k, v := range ssMap {
		L = append(L, GaugeValue("ss."+k, v))
	}

	return
}
