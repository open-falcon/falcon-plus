package funcs

import (
	"github.com/open-falcon/common/model"
	"github.com/toolkits/nux"
	"log"
)

func UdpMetrics() []*model.MetricValue {
	udp, err := nux.Snmp("Udp")
	if err != nil {
		log.Println("read snmp fail", err)
		return []*model.MetricValue{}
	}

	count := len(udp)
	ret := make([]*model.MetricValue, count)
	i := 0
	for key, val := range udp {
		ret[i] = CounterValue("snmp.Udp."+key, val)
		i++
	}

	return ret
}
