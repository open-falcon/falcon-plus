package cron

import (
	"fmt"
	"github.com/open-falcon/agent/g"
	"github.com/open-falcon/common/model"
	"log"
	"time"
)

func Report() {
	if g.Config().Heartbeat.Enabled && g.Config().Heartbeat.Addr != "" {
		go report(time.Duration(g.Config().Heartbeat.Interval) * time.Second)
	}
}

func report(interval time.Duration) {
	hostname, err := g.Hostname()
	if err != nil {
		hostname = fmt.Sprintf("error:%s", err.Error())
	}

	for {
		req := model.AgentReportRequest{
			Hostname:      hostname,
			IP:            g.IP(),
			AgentVersion:  g.VERSION,
			PluginVersion: g.GetCurrPluginVersion(),
		}

		var resp model.SimpleRpcResponse
		err = g.HbsClient.Call("Agent.ReportStatus", req, &resp)
		if err != nil || resp.Code != 0 {
			log.Println("call Agent.ReportStatus fail:", err, "Request:", req, "Response:", resp)
		}

		time.Sleep(interval)
	}
}
