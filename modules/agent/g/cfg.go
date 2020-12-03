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

package g

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/toolkits/file"
)

type PluginConfig struct {
	Enabled bool   `json:"enabled"`
	Dir     string `json:"dir"`
	Git     string `json:"git"`
	LogDir  string `json:"logs"`
}

type HeartbeatConfig struct {
	Enabled  bool   `json:"enabled"`
	Addr     string `json:"addr"`
	Interval int    `json:"interval"`
	Timeout  int    `json:"timeout"`
}

type TransferConfig struct {
	Enabled  bool     `json:"enabled"`
	Addrs    []string `json:"addrs"`
	Interval int      `json:"interval"`
	Timeout  int      `json:"timeout"`
}

type HttpConfig struct {
	Enabled  bool   `json:"enabled"`
	Listen   string `json:"listen"`
	Backdoor bool   `json:"backdoor"`
}

type CollectorConfig struct {
	IfacePrefix []string `json:"ifacePrefix"`
	MountPoint  []string `json:"mountPoint"`
}

type GlobalConfig struct {
	MaxMemory     uint64            `json:"max_memory"`
	MemoryCtrl    bool              `json:"memoryctrl"`
	Debug         bool              `json:"debug"`
	Hostname      string            `json:"hostname"`
	IP            string            `json:"ip"`
	Batch         int               `json:"batch,omitempty"`
	Plugin        *PluginConfig     `json:"plugin"`
	Heartbeat     *HeartbeatConfig  `json:"heartbeat"`
	Transfer      *TransferConfig   `json:"transfer"`
	Http          *HttpConfig       `json:"http"`
	Collector     *CollectorConfig  `json:"collector"`
	DefaultTags   map[string]string `json:"default_tags"`
	IgnoreMetrics map[string]bool   `json:"ignore"`
}

var (
	ConfigFile string
	config     *GlobalConfig
	lock       = new(sync.RWMutex)
)

func Config() *GlobalConfig {
	lock.RLock()
	defer lock.RUnlock()
	return config
}

func Hostname() (string, error) {
	hostname := Config().Hostname
	if hostname != "" {
		return hostname, nil
	}

	if os.Getenv("FALCON_ENDPOINT") != "" {
		hostname = os.Getenv("FALCON_ENDPOINT")
		return hostname, nil
	}

	hostname, err := os.Hostname()
	if err != nil {
		log.Println("ERROR: os.Hostname() fail", err)
	}
	return hostname, err
}

func IP() string {
	ip := Config().IP
	if ip != "" {
		// use ip in configuration
		return ip
	}

	if len(LocalIp) > 0 {
		ip = LocalIp
	}

	return ip
}

func ParseConfig(cfg string) {
	if cfg == "" {
		log.Fatalln("use -c to specify configuration file")
	}

	if !file.IsExist(cfg) {
		log.Fatalln("config file:", cfg, "is not existent. maybe you need `mv cfg.example.json cfg.json`")
	}

	ConfigFile = cfg

	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		log.Fatalln("read config file:", cfg, "fail:", err)
	}

	var c GlobalConfig
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		log.Fatalln("parse config file:", cfg, "fail:", err)
	}

	memCtrl :=os.Getenv("MEMORY_CTRL")
	if len(memCtrl)!=0{
		if strings.ToLower(memCtrl)=="true"{
			c.MemoryCtrl=true
		}else{
			c.MemoryCtrl=false
		}
	}

	transferAddr :=os.Getenv("TRANSFER_URL")
	if len(transferAddr)!=0{
		c.Transfer.Addrs=strings.Split(transferAddr,",")
		log.Println("transfer url: ",transferAddr)
	}

	heartbeatURL :=os.Getenv("HEARTBEAT_URL")
	if len(heartbeatURL)!=0{
		c.Heartbeat.Addr=heartbeatURL
		log.Println("heartbeat URL: ",transferAddr)
	}
	limitBatch :=os.Getenv("LIMIT_BATCH")
	if len(limitBatch)!=0{
		c.Batch,err=strconv.Atoi(limitBatch)
		if err!=nil{
			log.Println("invalid limit Batch: ",limitBatch)
		}
		log.Println("from env set Batch size: ", c.Batch)
	}
	if c.Batch<=0{
		c.Batch=2000
		log.Println("set batch default size: ",c.Batch)
	}

	lock.Lock()
	defer lock.Unlock()

	config = &c

	log.Println("read config file:", cfg, "successfully")
}
