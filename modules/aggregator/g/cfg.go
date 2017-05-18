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
	"github.com/toolkits/file"
	"log"
	"sync"
)

type HttpConfig struct {
	Enabled bool   `json:"enabled"`
	Listen  string `json:"listen"`
}

type DatabaseConfig struct {
	Addr     string `json:"addr"`
	Idle     int    `json:"idle"`
	Ids      []int  `json:"ids"`
	Interval int64  `json:"interval"`
}

type ApiConfig struct {
	ConnectTimeout int32  `json:"connect_timeout"`
	RequestTimeout int32  `json:"request_timeout"`
	PlusApi        string `json:"plus_api"`
	PlusApiToken   string `json:"plus_api_token"`
	PushApi        string `json:"push_api"`
}

type GlobalConfig struct {
	Debug    bool            `json:"debug"`
	Http     *HttpConfig     `json:"http"`
	Database *DatabaseConfig `json:"database"`
	Api      *ApiConfig      `json:"api"`
}

var (
	ConfigFile string
	config     *GlobalConfig
	configLock = new(sync.RWMutex)
)

func Config() *GlobalConfig {
	configLock.RLock()
	defer configLock.RUnlock()
	return config
}

func ParseConfig(cfg string) {
	if cfg == "" {
		log.Fatalln("use -c to specify configuration file")
	}

	if !file.IsExist(cfg) {
		log.Fatalln("config file:", cfg, "is not existent")
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

	configLock.Lock()
	defer configLock.Unlock()
	config = &c
	log.Println("read config file:", cfg, "successfully")
}
