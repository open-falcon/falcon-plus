package cron

import (
	"github.com/open-falcon/agent/g"
	"github.com/open-falcon/common/model"
	"log"
	"time"
)

func SyncTrustableIps() {
	if g.Config().Heartbeat.Enabled && g.Config().Heartbeat.Addr != "" {
		go syncTrustableIps()
	}
}

func syncTrustableIps() {

	duration := time.Duration(g.Config().Heartbeat.Interval) * time.Second

	for {
		time.Sleep(duration)

		var ips string
		err := g.HbsClient.Call("Agent.TrustableIps", model.NullRpcRequest{}, &ips)
		if err != nil {
			log.Println("ERROR: call Agent.TrustableIps fail", err)
			continue
		}

		g.SetTrustableIps(ips)
	}
}
