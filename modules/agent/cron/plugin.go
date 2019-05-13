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
	"github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/agent/g"
	"github.com/open-falcon/falcon-plus/modules/agent/plugins"
	"github.com/toolkits/file"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func SyncMinePlugins() {
	if !g.Config().Plugin.Enabled {
		return
	}

	// With hbs.
	if !g.Config().Plugin.Standalone {
		if !g.Config().Heartbeat.Enabled {
			return
		}

		if g.Config().Heartbeat.Addr == "" {
			return
		}
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

		if !g.Config().Plugin.Standalone {
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
				log.Println("call Agent.MinePlugin fail:", err)
				continue
			}

			if resp.Timestamp <= timestamp {
				continue
			}

			pluginDirs = resp.Plugins
			timestamp = resp.Timestamp

			if g.Config().Debug {
				log.Printf("call Agent.MinePlugin:%v\n", resp)
			}

			if len(pluginDirs) == 0 {
				plugins.ClearAllPlugins()
				continue
			}
		} else {
			pluginDirs = []string{"."}
		}

		desiredAll := make(map[string]*plugins.Plugin)
		filefmt_scripts := [][]string{}
		dirfmt_scripts := []string{}

		for _, script_path := range pluginDirs {
			//script_path could be a DIR or a SCRIPT_FILE_WITH_OR_WITHOUT_ARGS
			//比如： sys/ntp/60_ntp.py(arg1,arg2) 或者 sys/ntp/60_ntp.py 或者 sys/ntp
			//1. 参数只对单个脚本文件生效，目录不支持参数
			//2. 如果某个目录下的某个脚本被单独绑定到某个机器，那么再次绑定该目录时，该文件会不会再次执行
			var args string = ""

			re := regexp.MustCompile(`(.*)\((.*)\)`)
			path_args := re.FindAllStringSubmatch(script_path, -1)
			if path_args != nil {
				script_path = path_args[0][1]
				args = path_args[0][2]
			}

			abs_path := filepath.Join(g.Config().Plugin.Dir, script_path)
			if !file.IsExist(abs_path) {
				continue
			}

			if file.IsFile(abs_path) {
				filefmt_scripts = append(filefmt_scripts, []string{script_path, args})
				continue
			}

			dirfmt_scripts = append(dirfmt_scripts, script_path)
		}

		taken := make(map[string]struct{})
		for _, script_file := range filefmt_scripts {
			abs_path := filepath.Join(g.Config().Plugin.Dir, script_file[0])
			_, file_name := filepath.Split(abs_path)
			arr := strings.Split(file_name, "_")
			var cycle int
			var err error
			cycle, err = strconv.Atoi(arr[0])
			if err == nil {
				fi, _ := os.Stat(abs_path)
				plugin := &plugins.Plugin{FilePath: script_file[0], MTime: fi.ModTime().Unix(), Cycle: cycle, Args: script_file[1]}
				desiredAll[script_file[0]+"("+script_file[1]+")"] = plugin
			}
			//针对某个 hostgroup 绑定了单个脚本后，再绑定该脚本的目录时，会忽略目录中的该文件
			taken[script_file[0]] = struct{}{}
		}

		for _, script_path := range dirfmt_scripts {
			ps := plugins.ListPlugins(strings.Trim(script_path, "/"))
			for k, p := range ps {
				if _, ok := taken[k]; ok {
					continue
				}
				desiredAll[k] = p
			}
		}

		plugins.DelNoUsePlugins(desiredAll)
		plugins.AddNewPlugins(desiredAll)

		if g.Config().Debug {
			log.Printf("current plugins:%v\n", plugins.Plugins)
		}
	}
}
