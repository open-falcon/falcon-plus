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

type DBConfig struct {
	Dsn      string `json:"dsn"`
	MaxIdle  int    `json:"maxIdle"`
	Interval int    `json:"interval"`
}

type IndexConfig struct {
	Enabled bool `json:"enabled"`
}

type TransferConfig struct {
	Enabled bool     `json:"enabled"`
	Cluster []string `json:"cluster"`
}

type GlobalConfig struct {
	Debug    bool            `json:"debug"`
	Http     *HttpConfig     `json:"http"`
	DB       *DBConfig       `json:"db"`
	Index    *IndexConfig    `json:"index"`
	Transfer *TransferConfig `json:"transfer"`
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

	configLock.Lock()
	defer configLock.Unlock()
	config = &c

	log.Println("g:ParseConfig, ok, ", cfg)
}
