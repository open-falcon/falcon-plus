package cron

import (
	"fmt"
	"github.com/open-falcon/agent/g"
	"time"
)

func Report() {
	if g.Config().Heartbeat.Enabled && g.Config().Heartbeat.Addr != "" {
		go report(time.Duration(g.Config().Heartbeat.Interval) * time.Second)
	}
}

func report(interval time.Duration) {
	ip := ""
	if len(g.LocalIps) > 0 {
		ip = g.LocalIps[0]
	}

	hostname, err := g.Hostname()
	if err != nil {
		hostname = fmt.Sprintf("error:%s", err.Error())
	}

	for {
		req := g.AgentReportReq{
			HostName:      hostname,
			Version:       g.VERSION,
			Meta:          ip,
			PluginVersion: GetCurrPluginVersion(),
		}

		var resp g.AgentReportResp
		g.HbsClient.Call("Agent.ReportStatus", req, &resp)

		time.Sleep(interval)
	}
}
