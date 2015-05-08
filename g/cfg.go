package g

import (
	"encoding/json"
	"github.com/toolkits/file"
	"github.com/toolkits/logger"
	"log"
	"sync"
)

type HttpConfig struct {
	Enabled bool   `json:"enabled"`
	Listen  string `json:"listen"`
}

type GraphConfig struct {
	Backends       string `json:"backends"`
	ReloadInterval int    `json:"reload_interval"`
	Timeout        int    `json:"timeout"`   // millisecond, connect timeout or request timeout
	MaxConns       int    `json:"max_conns"` // 链接池
	MaxIdle        int    `json:"max_idle"`
	Replicas       int    `json:"replicas"`
}

type GlobalConfig struct {
	LogLevel string       `json:"log_level"`
	SlowLog  int          `json:"slowlog"` // 单位ms，耗时超过这个的所有转发会被记录到日志
	Graph    *GraphConfig `json:"graph"`
	Http     *HttpConfig  `json:"http"`
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

	logger.SetLevel(config.LogLevel)
	log.Println("read config file:", cfg, "successfully")
}
