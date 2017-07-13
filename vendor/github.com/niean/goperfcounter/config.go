package goperfcounter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

type GlobalConfig struct {
	Debug    bool        `json:"debug"`
	Hostname string      `json:"hostname"`
	Tags     string      `json:"tags"`
	Step     int64       `json:"step"`
	Bases    []string    `json:"bases"`
	Push     *PushConfig `json:"push"`
	Http     *HttpConfig `json:"http"`
}
type HttpConfig struct {
	Enabled bool   `json:"enabled"`
	Listen  string `json:"listen"`
}
type PushConfig struct {
	Enabled bool   `json:"enabled"`
	Api     string `json:"api"`
}

var (
	configFn     = "./perfcounter.json"
	defaultTags  = ""
	defaultStep  = int64(60) //time in sec
	defaultBases = []string{}
	defaultPush  = &PushConfig{Enabled: true, Api: "http://127.0.0.1:1988/v1/push"}
	defaultHttp  = &HttpConfig{Enabled: false, Listen: ""}
)

var (
	cfg     *GlobalConfig
	cfgLock = new(sync.RWMutex)
)

//
func config() *GlobalConfig {
	cfgLock.RLock()
	defer cfgLock.RUnlock()
	return cfg
}

func loadConfig() error {
	if !isFileExist(configFn) {
		return fmt.Errorf("config file not found: %s", configFn)
	}

	c, err := parseConfig(configFn)
	if err != nil {
		return err
	}

	updateConfig(c)
	return nil
}

func setDefaultConfig() {
	dcfg := defaultConfig()
	updateConfig(dcfg)
}

func defaultConfig() GlobalConfig {
	return GlobalConfig{
		Debug:    false,
		Hostname: defaultHostname(),
		Tags:     defaultTags,
		Step:     defaultStep,
		Bases:    defaultBases,
		Push:     defaultPush,
		Http:     defaultHttp,
	}
}

//
func updateConfig(c GlobalConfig) {
	nc := formatConfig(c)
	cfgLock.Lock()
	defer cfgLock.Unlock()
	cfg = &nc
}

func formatConfig(c GlobalConfig) GlobalConfig {
	nc := c
	if nc.Hostname == "" {
		nc.Hostname = defaultHostname()
	}
	if nc.Step < 1 {
		nc.Step = defaultStep
	}
	if nc.Tags != "" {
		tagsOk := true
		tagsSlice := strings.Split(nc.Tags, ",")
		for _, tag := range tagsSlice {
			kv := strings.Split(tag, "=")
			if len(kv) != 2 || kv[0] == "name" { // name是保留tag
				tagsOk = false
				break
			}
		}
		if !tagsOk {
			nc.Tags = defaultTags
		}
	}
	if nc.Push.Enabled && nc.Push.Api == "" {
		nc.Push = defaultPush
	}
	if len(nc.Bases) < 1 {
		nc.Bases = defaultBases
	}

	return nc
}

func parseConfig(cfg string) (GlobalConfig, error) {
	var c GlobalConfig

	if cfg == "" {
		return c, fmt.Errorf("config file not found")
	}

	configContent, err := readFileString(cfg)
	if err != nil {
		return c, fmt.Errorf("read config file %s error: %v", cfg, err.Error())
	}

	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		return c, fmt.Errorf("parse config file %s error: %v", cfg, err.Error())
	}
	return c, nil
}

func defaultHostname() string {
	hostname, _ := os.Hostname()
	return hostname
}

func isFileExist(fn string) bool {
	_, err := os.Stat(fn)
	return err == nil || os.IsExist(err)
}

func readFileString(fn string) (string, error) {
	b, err := ioutil.ReadFile(fn)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(b)), nil
}
