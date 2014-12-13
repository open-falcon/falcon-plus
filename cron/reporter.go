package cron

import (
	"bytes"
	"fmt"
	"github.com/open-falcon/agent/g"
	"github.com/toolkits/file"
	"os"
	"os/exec"
	"strings"
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

	hostname, err := os.Hostname()
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

func GetCurrPluginVersion() string {
	if !g.Config().Plugin.Enabled {
		return "plugin not enabled"
	}

	pluginDir := g.Config().Plugin.Dir
	if !file.IsExist(pluginDir) {
		return "plugin dir not existent"
	}

	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = pluginDir

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return fmt.Sprintf("Error:%s", err.Error())
	}

	return strings.TrimSpace(out.String())
}
