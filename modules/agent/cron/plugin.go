// Copyright 2017 Xiaomi, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cron

import (
	"log"
	"strings"
	"time"

	"github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/agent/g"
	"github.com/open-falcon/falcon-plus/modules/agent/plugins"
)

func SyncMinePlugins() {
	if !g.Config().Plugin.Enabled {
		return
	}

	if !g.Config().Heartbeat.Enabled {
		return
	}

	if len(g.Config().Heartbeat.Addrs) == 0 {
		return
	}

	go syncMinePlugins()
}

func syncMinePlugins() {

	var (
		timestamp  int64 = -1
		pluginDirs []string
	)

	duration := time.Duration(g.Config().Heartbeat.Interval) * time.Second

	for {
		time.Sleep(duration)

		hostname, err := g.Hostname()
		if err != nil {
			continue
		}

		req := model.AgentHeartbeatRequest{
			Hostname: hostname,
		}

		var resp model.AgentPluginsResponse
		err = g.HbsClient.Call("Agent.MinePlugins", req, &resp)
		if err != nil {
			log.Println("ERROR:", err)
			continue
		}

		if resp.Timestamp <= timestamp {
			continue
		}

		pluginDirs = resp.Plugins
		timestamp = resp.Timestamp

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
