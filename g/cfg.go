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

type SocketConfig struct {
	Enabled bool   `json:"enabled"`
	Listen  string `json:"listen"`
	Timeout int    `json:"timeout"`
}

type JudgeConfig struct {
	Enabled     bool              `json:"enabled"`
	Batch       int               `json:"batch"`
	ConnTimeout int               `json:"connTimeout"`
	CallTimeout int               `json:"callTimeout"`
	PingMethod  string            `json:"pingMethod"`
	MaxConns    int               `json:"maxConns"`
	MaxIdle     int               `json:"maxIdle"`
	Replicas    int               `json:"replicas"`
	Cluster     map[string]string `json:"cluster"`
}

type GraphConfig struct {
	Enabled          bool              `json:"enabled"`
	Batch            int               `json:"batch"`
	ConnTimeout      int               `json:"connTimeout"`
	CallTimeout      int               `json:"callTimeout"`
	PingMethod       string            `json:"pingMethod"`
	MaxConns         int               `json:"maxConns"`
	MaxIdle          int               `json:"maxIdle"`
	Replicas         int               `json:"replicas"`
	Migrating        bool              `json:"migrating"`
	Cluster          map[string]string `json:"cluster"`
	ClusterMigrating map[string]string `json:"clusterMigrating"`
}

type GlobalConfig struct {
	Debug  bool          `json:"debug"`
	Http   *HttpConfig   `json:"http"`
	Rpc    *RpcConfig    `json:"rpc"`
	Socket *SocketConfig `json:"socket"`
	Judge  *JudgeConfig  `json:"judge"`
	Graph  *GraphConfig  `json:"graph"`
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
