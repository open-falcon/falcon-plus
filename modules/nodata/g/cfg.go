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

type PlusAPIConfig struct {
	Addr           string `json:"addr"`
	Token          string `json:"token"`
	ConnectTimeout int32  `json:"connectTimeout"`
	RequestTimeout int32  `json:"requestTimeout"`
}

type NdConfig struct {
	Enabled bool   `json:"enabled"`
	Dsn     string `json:"dsn"`
	MaxIdle int32  `json:"maxIdle"`
}

type CollectorConfig struct {
	Enabled    bool  `json:"enabled"`
	Batch      int32 `json:"batch"`
	Concurrent int32 `json:"concurrent"`
}

type SenderConfig struct {
	Enabled        bool   `json:"enabled"`
	TransferAddr   string `json:"transferAddr"`
	ConnectTimeout int32  `json:"connectTimeout"`
	RequestTimeout int32  `json:"requestTimeout"`
	Batch          int32  `json:"batch"`
	DelayPeriod      int32  `json:"delayPeriod"`
}

type GlobalConfig struct {
	Debug     bool             `json:"debug"`
	Http      *HttpConfig      `json:"http"`
	PlusApi   *PlusAPIConfig   `json:"plus_api"`
	Config    *NdConfig        `json:"config"`
	Collector *CollectorConfig `json:"collector"`
	Sender    *SenderConfig    `json:"sender"`
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
		log.Fatalln("g.ParseConfig error, parse config file", cfg, "fail,", err)
	}

	configLock.Lock()
	defer configLock.Unlock()
	config = &c

	log.Println("g.ParseConfig ok, file ", cfg)
}
