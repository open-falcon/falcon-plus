package g

import (
	"encoding/json"
	"github.com/toolkits/file"
	"log"
	"sync"
)

type PluginConfig struct {
	Dir    string `json:"dir"`
	Git    string `json:"git"`
	LogDir string `json:"logs"`
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
	Port int `json:"port"`
}

type GlobalConfig struct {
	Debug     bool             `json:"debug"`
	Hostname  string           `json:"hostname"`
	Plugin    *PluginConfig    `json:"plugin"`
	Heartbeat *HeartbeatConfig `json:"heartbeat"`
	Transfer  *TransferConfig  `json:"transfer"`
	Http      *HttpConfig      `json:"http"`
}

var (
	config *GlobalConfig
	lock   = new(sync.RWMutex)
)

func GetConfig() *GlobalConfig {
	lock.RLock()
	defer lock.RUnlock()
	return config
}

func ParseConfig(cfg string) {
	lock.Lock()
	defer lock.Unlock()
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

	config = &c

	if config.Debug {
		log.Println("read config file:", cfg, "successfully")
	}
}
