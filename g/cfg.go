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

type IndexConfig struct {
	Enabled bool     `json:"enabled"`
	Dsn     string   `json:"dsn"`
	MaxIdle int      `json:"maxIdle"`
	Cluster []string `json:"cluster"`
}

type MonitorConfig struct {
	Enabled bool     `json:"enabled"`
	MailUrl string   `json:"mailUrl"`
	MailTos string   `json:"mailTos"`
	Cluster []string `json:"cluster"`
}

type CollectorConfig struct {
	Enabled   bool     `json:"enabled"`
	DestUrl   string   `json:"destUrl"`
	SrcUrlFmt string   `json:"srcUrlFmt"`
	Cluster   []string `json:"cluster"`
}

type GlobalConfig struct {
	Debug     bool             `json:"debug"`
	Http      *HttpConfig      `json:"http"`
	Index     *IndexConfig     `json:"index"`
	Monitor   *MonitorConfig   `json:"monitor"`
	Collector *CollectorConfig `json:"collector"`
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
