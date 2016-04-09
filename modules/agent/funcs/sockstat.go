package funcs

import (
	"github.com/open-falcon/common/model"
	"github.com/toolkits/nux"
	"log"
)

func SocketStatSummaryMetrics() (L []*model.MetricValue) {
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
