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

package plugins

import (
	"bytes"
	"encoding/json"
	"log"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"

	"github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/agent/g"
	"github.com/toolkits/file"
	"github.com/toolkits/sys"
)

type PluginScheduler struct {
	Ticker *time.Ticker
	Plugin *Plugin
	Quit   chan struct{}
}

func NewPluginScheduler(p *Plugin) *PluginScheduler {
	scheduler := PluginScheduler{Plugin: p}
	scheduler.Ticker = time.NewTicker(time.Duration(p.Cycle) * time.Second)
	scheduler.Quit = make(chan struct{})
	return &scheduler
}

func (this *PluginScheduler) Schedule() {
	go func() {
		for {
			select {
			case <-this.Ticker.C:
				PluginRun(this.Plugin)
			case <-this.Quit:
				this.Ticker.Stop()
				return
			}
		}
	}()
}

func (this *PluginScheduler) Stop() {
	close(this.Quit)
}

func PluginRun(plugin *Plugin) {

	timeout := plugin.Cycle*1000 - 500
	fpath := filepath.Join(g.Config().Plugin.Dir, plugin.FilePath)

	if !file.IsExist(fpath) {
		log.Println("no such plugin:", fpath)
		return
	}

	debug := g.Config().Debug
	if debug {
		log.Println(fpath, "running...")
	}

	cmd := exec.Command(fpath)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	err := cmd.Start()
	if err != nil {
		log.Printf("[ERROR] plugin start fail, error: %s\n", err)
		return
	}
	if debug {
		log.Println("plugin started:", fpath)
	}

	err, isTimeout := sys.CmdRunWithTimeout(cmd, time.Duration(timeout)*time.Millisecond)

	errStr := stderr.String()
	if errStr != "" {
		logFile := filepath.Join(g.Config().Plugin.LogDir, plugin.FilePath+".stderr.log")
		if _, err = file.WriteString(logFile, errStr); err != nil {
			log.Printf("[ERROR] write log to %s fail, error: %s\n", logFile, err)
		}
	}

	if isTimeout {
		// has be killed
		if err == nil && debug {
			log.Println("[INFO] timeout and kill process", fpath, "successfully")
		}

		if err != nil {
			log.Println("[ERROR] kill process", fpath, "occur error:", err)
		}

		return
	}

	if err != nil {
		log.Println("[ERROR] exec plugin", fpath, "fail. error:", err)
		return
	}

	// exec successfully
	data := stdout.Bytes()
	if len(data) == 0 {
		if debug {
			log.Println("[DEBUG] stdout of", fpath, "is blank")
		}
		return
	}

	var metrics []*model.MetricValue
	err = json.Unmarshal(data, &metrics)
	if err != nil {
		log.Printf("[ERROR] json.Unmarshal stdout of %s fail. error:%s stdout: \n%s\n", fpath, err, stdout.String())
		return
	}

	g.SendToTransfer(metrics)
}
