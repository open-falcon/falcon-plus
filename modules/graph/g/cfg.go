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
	"strconv"
	"sync/atomic"
	"unsafe"

	"github.com/toolkits/file"
)

type File struct {
	Filename string
	Body     []byte
}

type HttpConfig struct {
	Enabled bool   `json:"enabled"`
	Listen  string `json:"listen"`
}

type RpcConfig struct {
	Enabled bool   `json:"enabled"`
	Listen  string `json:"listen"`
}

type RRDConfig struct {
	Storage string `json:"storage"`
}

type DBConfig struct {
	Dsn     string `json:"dsn"`
	MaxIdle int    `json:"maxIdle"`
}

type GlobalConfig struct {
	Pid            string      `json:"pid"`
	Debug          bool        `json:"debug"`
	Http           *HttpConfig `json:"http"`
	Rpc            *RpcConfig  `json:"rpc"`
	RRD            *RRDConfig  `json:"rrd"`
	DB             *DBConfig   `json:"db"`
	CallTimeout    int32       `json:"callTimeout"`
	IOWorkerNum    int         `json:"ioWorkerNum"`
	FirstBytesSize int
	Migrate        struct {
		Concurrency int               `json:"concurrency"` //number of multiple worker per node
		Enabled     bool              `json:"enabled"`
		Replicas    int               `json:"replicas"`
		Cluster     map[string]string `json:"cluster"`
	} `json:"migrate"`
}

var (
	ConfigFile string
	ptr        unsafe.Pointer
)

func Config() *GlobalConfig {
	return (*GlobalConfig)(atomic.LoadPointer(&ptr))
}

func ParseConfig(cfg string) {
	if cfg == "" {
		log.Fatalln("config file not specified: use -c $filename")
	}

	if !file.IsExist(cfg) {
		log.Fatalln("config file specified not found:", cfg)
	}

	ConfigFile = cfg

	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		log.Fatalln("read config file", cfg, "error:", err.Error())
	}

	var c GlobalConfig
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		log.Fatalln("parse config file", cfg, "error:", err.Error())
	}

	if c.Migrate.Enabled && len(c.Migrate.Cluster) == 0 {
		c.Migrate.Enabled = false
	}

	// 确保ioWorkerNum是2^N
	if c.IOWorkerNum == 0 || (c.IOWorkerNum&(c.IOWorkerNum-1) != 0) {
		log.Fatalf("IOWorkerNum must be 2^N, current IOWorkerNum is %v", c.IOWorkerNum)
	}

	// 需要md5的前多少位参与ioWorker的分片计算
	c.FirstBytesSize = len(strconv.FormatInt(int64(c.IOWorkerNum), 16))

	// set config
	atomic.StorePointer(&ptr, unsafe.Pointer(&c))

	log.Println("g.ParseConfig ok, file", cfg)
}
