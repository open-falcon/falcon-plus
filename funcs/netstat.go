package funcs

import (
	"github.com/open-falcon/agent/g"
	"github.com/toolkits/nux"
	"log"
)

var USES = map[string]struct{}{
	"PruneCalled":        struct{}{},
	"LockDroppedIcmps":   struct{}{},
	"ArpFilter":          struct{}{},
	"TW":                 struct{}{},
	"DelayedACKLocked":   struct{}{},
	"ListenOverflows":    struct{}{},
	"ListenDrops":        struct{}{},
	"TCPPrequeueDropped": struct{}{},
	"TCPTSReorder":       struct{}{},
	"TCPDSACKUndo":       struct{}{},
	"TCPLoss":            struct{}{},
	"TCPLostRetransmit":  struct{}{},
	"TCPLossFailures":    struct{}{},
	"TCPFastRetrans":     struct{}{},
	"TCPTimeouts":        struct{}{},
	"TCPSchedulerFailed": struct{}{},
	"TCPAbortOnMemory":   struct{}{},
	"TCPAbortOnTimeout":  struct{}{},
	"TCPAbortFailed":     struct{}{},
	"TCPMemoryPressures": struct{}{},
	"TCPSpuriousRTOs":    struct{}{},
	"TCPBacklogDrop":     struct{}{},
	"TCPMinTTLDrop":      struct{}{},
}

func NetstatMetrics() []*g.MetricValue {
	tcpExts, err := nux.Netstat("TcpExt")

	ret := make([]*g.MetricValue, 0)

	if err != nil {
		log.Println(err)
		return ret
	}

	cnt := len(tcpExts)
	if cnt == 0 {
		return ret
	}

	for key, val := range tcpExts {
		if _, ok := USES[key]; !ok {
			continue
		}
		ret = append(ret, CounterValue("TcpExt."+key, val))
	}

	return ret
}
