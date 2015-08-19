package g

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/toolkits/file"
)

type HttpConfig struct {
	Enable bool   `json:"enable"`
	Listen string `json:"listen"`
}

type RpcConfig struct {
	Enable bool   `json:"enable"`
	Listen string `json:"listen"`
}

type SocketConfig struct {
	Enable  bool   `json:"enable"`
	Listen  string `json:"listen"`
	Timeout int32  `json:"timeout"`
}

type TransferConfig struct {
	Enable      bool   `json:"enable"`
	Batch       int32  `json:"batch"`
	ConnTimeout int32  `json:"connTimeout"`
	CallTimeout int32  `json:"callTimeout"`
	MaxConns    int32  `json:"maxConns"`
	MaxIdle     int32  `json:"maxIdle"`
	Addr        string `json:"addr"`
}

type GlobalConfig struct {
	Debug    bool            `json:"debug"`
	Http     *HttpConfig     `json:"http"`
	Rpc      *RpcConfig      `json:"rpc"`
	Socket   *SocketConfig   `json:"socket"`
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

	// 配置文件正确性 校验, 不合法则直接 Exit(1)
	// TODO

	configLock.Lock()
	defer configLock.Unlock()
	config = &c

	log.Println("g.ParseConfig ok, file ", cfg)
}
