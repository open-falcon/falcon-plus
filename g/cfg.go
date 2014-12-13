package g

import (
	"encoding/json"
	"github.com/toolkits/file"
	"log"
	"sync"
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
	Enabled  bool   `json:"enabled"`
	Addr     string `json:"addr"`
	Interval int    `json:"interval"`
	Timeout  int    `json:"timeout"`
}

type HttpConfig struct {
	Enabled bool `json:"enabled"`
	Port    int  `json:"port"`
}

type CollectorConfig struct {
	IfacePrefix []string `json:"ifacePrefix"`
}

type GlobalConfig struct {
	Debug     bool             `json:"debug"`
	Hostname  string           `json:"hostname"`
	Plugin    *PluginConfig    `json:"plugin"`
	Heartbeat *HeartbeatConfig `json:"heartbeat"`
	Transfer  *TransferConfig  `json:"transfer"`
	Http      *HttpConfig      `json:"http"`
	Collector *CollectorConfig `json:"collector"`
}

var (
	config *GlobalConfig
	lock   = new(sync.RWMutex)
)

func Config() *GlobalConfig {
	lock.RLock()
	defer lock.RUnlock()
	return config
}

func ParseConfig(cfg string) {
	if cfg == "" {
		log.Fatalln("use -c to specify configuration file")
	}

	if !file.IsExist(cfg) {
		log.Fatalln("config file:", cfg, "is not existent")
	}

	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		log.Fatalln("read config file:", cfg, "fail:", err)
	}

	var c GlobalConfig
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		log.Fatalln("parse config file:", cfg, "fail:", err)
	}

	lock.Lock()
	defer lock.Unlock()

	config = &c

	if config.Debug {
		log.Println("read config file:", cfg, "successfully")
	}
}
