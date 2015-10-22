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

type QueryConfig struct {
	QueryAddr      string `json:"queryAddr"`
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

type BlockConfig struct {
	Enabled        bool    `json:"enabled"`
	Threshold      int32   `json:"threshold"`
	SetBlock       bool    `json:"setBlock"`
	EnableGauss    bool    `json:"enableGauss"`
	Hostname       string  `json:"hostname"`
	FloodCounter   string  `json:"floodCounter"`
	GaussFilter    float64 `json:"gaussFilter"`
	Gauss3SigmaMin float64 `json:"gauss3SigmaMin"`
	Gauss3SigmaMax float64 `json:"gauss3SigmaMax"`
}

type SenderConfig struct {
	Enabled        bool         `json:"enabled"`
	TransferAddr   string       `json:"transferAddr"`
	ConnectTimeout int32        `json:"connectTimeout"`
	RequestTimeout int32        `json:"requestTimeout"`
	Batch          int32        `json:"batch"`
	Block          *BlockConfig `json:"block"`
}

type GlobalConfig struct {
	Debug     bool             `json:"debug"`
	Http      *HttpConfig      `json:"http"`
	Query     *QueryConfig     `json:"query"`
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
