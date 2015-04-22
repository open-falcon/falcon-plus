package cron

import (
	"github.com/open-falcon/agent/g"
	"github.com/open-falcon/agent/plugins"
	"github.com/open-falcon/common/model"
	"log"
	"strings"
	"time"
)

func SyncPlugin() {
	if !g.Config().Plugin.Enabled {
		return
	}

	if !g.Config().Heartbeat.Enabled {
		return
	}

	if g.Config().Heartbeat.Addr == "" {
		return
	}

	go syncPlugin()
}

func syncPlugin() {

	var (
		checksum   string = "nil"
		timestamp  int64  = -1
		pluginDirs []string
	)

	duration := time.Duration(g.Config().Heartbeat.Interval) * time.Second

	for {
	REST:
		time.Sleep(duration)

		hostname, err := g.Hostname()
		if err != nil {
			goto REST
		}

		req := model.AgentHeartbeatRequest{
			Hostname: hostname,
			Checksum: checksum,
		}

		var resp model.AgentPluginsResponse
		err = g.HbsClient.Call("Agent.MinePlugins", req, &resp)
		if err != nil {
			log.Println("ERROR:", err)
			goto REST
		}

		if resp.Timestamp <= timestamp {
			goto REST
		}

		if resp.Checksum == checksum {
			goto REST
		}

		pluginDirs = resp.Plugins
		timestamp = resp.Timestamp
		checksum = resp.Checksum

		if g.Config().Debug {
			log.Println(&resp)
		}

		if len(pluginDirs) == 0 {
			plugins.ClearAllPlugins()
		}

		desiredAll := make(map[string]*plugins.Plugin)

		for _, p := range pluginDirs {
			underOneDir := plugins.ListPlugins(strings.Trim(p, "/"))
			for k, v := range underOneDir {
				desiredAll[k] = v
			}
		}

		plugins.DelNoUsePlugins(desiredAll)
		plugins.AddNewPlugins(desiredAll)

	}
}
