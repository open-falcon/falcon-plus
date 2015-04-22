package cron

import (
	"github.com/open-falcon/agent/g"
	"github.com/open-falcon/agent/plugins"
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
		localCheckSum  string
		localTimestamp int64
		pluginPaths    []*g.Plugin
	)

	duration := time.Duration(g.Config().Heartbeat.Interval) * time.Second

	for {
	REST:
		time.Sleep(duration)

		hostname, err := g.Hostname()
		if err != nil {
			log.Println("[ERROR]", err)
			goto REST
		}

		req := g.AgentReq{
			Host:     g.Host{HostName: hostname},
			Checksum: localCheckSum,
		}

		var resp g.AgentPluginsResp
		err = g.HbsClient.Call("Agent.GetPlugins", req, &resp)
		if err != nil {
			log.Println("[ERROR]", err)
			goto REST
		}

		if resp.Checksum == "" {
			log.Println("[ERROR] resp.Checksum is blank")
			goto REST
		}

		if resp.Timestamp <= localTimestamp {
			log.Println("[ERROR] resp.Timestamp <= localTimestamp")
			goto REST
		}

		if resp.Checksum == localCheckSum {
			goto REST
		}

		pluginPaths = resp.Plugins
		localTimestamp = resp.Timestamp
		localCheckSum = resp.Checksum

		if g.Config().Debug {
			log.Println("Plugins::::::::::::::::::::")
			log.Println("PluginPaths:", pluginPaths)
			log.Println("Timestamp:", localTimestamp)
			log.Println("Checksum:", localCheckSum)
		}

		if len(pluginPaths) == 0 {
			plugins.ClearAllPlugins()
		}

		desiredAll := make(map[string]*plugins.Plugin)

		pluginVersion := g.GetCurrPluginVersion()

		for _, p := range pluginPaths {
			underOneDir := plugins.PluginsUnder(strings.Trim(p.Path, "/"), pluginVersion)
			for k, v := range underOneDir {
				desiredAll[k] = v
			}
		}

		plugins.DelNoUsePlugins(desiredAll)
		plugins.AddNewPlugins(desiredAll)

	}
}
