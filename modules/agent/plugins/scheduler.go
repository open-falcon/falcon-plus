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

	"fmt"
	"github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/agent/g"
	"github.com/toolkits/file"
	"github.com/toolkits/sys"
	"strconv"
	"strings"
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

	metrics := processOutput(plugin, data)
	g.SendToTransfer(metrics)
}

func processOutput(plugin *Plugin, data []byte) []*model.MetricValue {
	var metrics []*model.MetricValue

	var err error
	dataStr := string(data)
	// 尝试多种方式对输出进行解析
	if strings.TrimSpace(dataStr)[0] == '[' {
		err = parseAsJSON(data, &metrics)
	} else {
		var metric *model.MetricValue
		for _, line := range strings.Split(string(data), "\n") {
			line = strings.TrimSpace(line)
			if len(line) == 0 {
				continue
			}

			if line[0] == '{' {
				metric, err = parseAsLineJSON(line)
			} else {
				metric, err = parseAsLineValue(line)
			}

			if err == nil {
				metrics = append(metrics, metric)
			} else {
				log.Printf("[ERROR] %v, output line: \"%s\"", err, line)
			}
		}
	}

	// 没有任何数据成功解析
	if len(metrics) == 0 {
		log.Println("[ERROR] Invalid output:\n", string(data))
	}

	// 插件的 step 需要在这里设置,
	// 其他的如: Timestamp, Endpoint, Type 会在 SendToTransfer() 中设置
	for _, m := range metrics {
		if m.Step == 0 {
			m.Step = int64(plugin.Cycle)
		}
	}

	return metrics
}

/*将插件的整个输出当作 json 处理*/
func parseAsJSON(data []byte, metrics *[]*model.MetricValue) error {
	err := json.Unmarshal(data, metrics)
	if err != nil {
		if g.Config().Debug {
			log.Printf("[DEBUG] Try parseAsJSON(): %s, IGNORE\n", err)
		}
		return err
	}
	return nil
}

/*将插件的一行输出当作 json 处理*/
func parseAsLineJSON(line string) (*model.MetricValue, error) {
	var metric model.MetricValue
	err := json.Unmarshal([]byte(line), &metric)
	return &metric, err
}

/*
* 将插件的一行输出当作一个数值处理
* 行格式 metric tag1=v1,tag2=v2 value timestamp
* 可选格式 :
* 1. metric value
* 2. metric tag1=v1,tag2=v2 value
* 3. metric tag1=v1,tag2=v2 value timestamp
* 4. metric value timestamp
 */
func parseAsLineValue(line string) (*model.MetricValue, error) {
	metric := model.MetricValue{Timestamp: time.Now().Unix()}
	fields := strings.Fields(line)
	timestampStr := fmt.Sprintf("%d", metric.Timestamp)

	if g.Config().Debug {
		log.Printf("[DEBUG] parseAsLineValue: %s\n", line)
	}

	if len(fields) == 2 { // 格式 1
		metric.Metric = fields[0]
		metric.Value = fields[1]
	} else if len(fields) == 3 {
		if strings.Index(fields[1], "=") > 0 {
			// 格式 2
			metric.Metric = fields[0]
			metric.Tags = fields[1]
			metric.Value = fields[2]
		} else {
			// 格式 4
			metric.Metric = fields[0]
			metric.Value = fields[1]
			timestampStr = fields[2]
		}
	} else if len(fields) == 4 { // 格式 3
		metric.Metric = fields[0]
		metric.Tags = fields[1]
		metric.Value = fields[2]
		timestampStr = fields[3]
	} else {
		return nil, fmt.Errorf("invalid line")
	}

	if val, err := strconv.Atoi(metric.Value.(string)); err == nil {
		metric.Value = val
	} else {
		return nil, fmt.Errorf("invalid metric value")
	}

	if timestamp, err := strconv.ParseInt(timestampStr, 10, 64); err == nil {
		metric.Timestamp = timestamp
	} else {
		return nil, fmt.Errorf("invalid metric timestamp")
	}

	return &metric, nil
}
