package funcs

import (
	"fmt"
	"github.com/open-falcon/agent/g"
	"github.com/open-falcon/common/model"
	"github.com/toolkits/nux"
	"github.com/toolkits/slice"
	"log"
)

func PortMetrics() (L []*model.MetricValue) {

	reportPorts := g.ReportPorts()
	sz := len(reportPorts)
	if sz == 0 {
		return
	}

	allListeningPorts, err := nux.ListeningPorts()
	if err != nil {
		log.Println(err)
		return
	}

	for i := 0; i < sz; i++ {
		tags := fmt.Sprintf("port=%d", reportPorts[i])
		if slice.ContainsInt64(allListeningPorts, reportPorts[i]) {
			L = append(L, GaugeValue("net.port.listen", 1, tags))
		} else {
			L = append(L, GaugeValue("net.port.listen", 0, tags))
		}
	}

	return
}
