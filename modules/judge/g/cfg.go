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

type RpcConfig struct {
	Enabled bool   `json:"enabled"`
	Listen  string `json:"listen"`
}

type HbsConfig struct {
	Servers     []string `json:"servers"`
	Timeout     int64    `json:"timeout"`
	Interval    int64    `json:"interval"`
	CallTimeout int64    `json:"callTimeout"`
}

type RedisConfig struct {
	Dsn          string `json:"dsn"`
	MaxIdle      int    `json:"maxIdle"`
	ConnTimeout  int    `json:"connTimeout"`
	ReadTimeout  int    `json:"readTimeout"`
	WriteTimeout int    `json:"writeTimeout"`
}

type AlarmConfig struct {
	Enabled      bool         `json:"enabled"`
	MinInterval  int64        `json:"minInterval"`
	QueuePattern string       `json:"queuePattern"`
	Redis        *RedisConfig `json:"redis"`
}

type GlobalConfig struct {
	Debug     bool         `json:"debug"`
	DebugHost string       `json:"debugHost"`
	Remain    int          `json:"remain"`
	Http      *HttpConfig  `json:"http"`
	Rpc       *RpcConfig   `json:"rpc"`
	Hbs       *HbsConfig   `json:"hbs"`
	Alarm     *AlarmConfig `json:"alarm"`
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
