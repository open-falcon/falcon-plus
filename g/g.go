package g

import (
//	"io/ioutil"
)

var Modules map[string]bool
var ModuleBins map[string]string
var ModuleConfs map[string]string
var ModuleApps map[string]string
var AllModulesInOrder []string

func init() {
	//	dirs, _ := ioutil.ReadDir("./modules")

	//	for _, dir := range dirs {
	//		Modules[dir.Name()] = true
	//	}
	Modules = map[string]bool{
		"agent":      true,
		"aggregator": true,
		"fe":         true,
		"graph":      true,
		"hbs":        true,
		"judge":      true,
		"nodata":     true,
		"query":      true,
		"sender":     true,
		"task":       true,
		"transfer":   true,
		"gateway":    true,
	}

	ModuleBins = map[string]string{
		"agent":      "./agent/bin/falcon-agent",
		"aggregator": "./aggregator/bin/falcon-aggregator",
		"fe":         "./fe/bin/falcon-fe",
		"graph":      "./graph/bin/falcon-graph",
		"hbs":        "./hbs/bin/falcon-hbs",
		"judge":      "./judge/bin/falcon-judge",
		"nodata":     "./nodata/bin/falcon-nodata",
		"query":      "./query/bin/falcon-query",
		"sender":     "./sender/bin/falcon-sender",
		"task":       "./task/bin/falcon-task",
		"transfer":   "./transfer/bin/falcon-transfer",
		"gateway":    "./gateway/bin/falcon-gateway",
	}

	ModuleConfs = map[string]string{
		"agent":      "./agent/config/cfg.json",
		"aggregator": "./aggregator/config/cfg.json",
		"fe":         "./fe/config/cfg.json",
		"graph":      "./graph/config/cfg.json",
		"hbs":        "./hbs/config/cfg.json",
		"judge":      "./judge/config/cfg.json",
		"nodata":     "./nodata/config/cfg.json",
		"query":      "./query/config/cfg.json",
		"sender":     "./sender/config/cfg.json",
		"task":       "./task/config/cfg.json",
		"transfer":   "./transfer/config/cfg.json",
		"gateway":    "./gateway/config/cfg.json",
	}

	ModuleApps = map[string]string{
		"agent":      "falcon-agent",
		"aggregator": "falcon-aggregator",
		"graph":      "falcon-graph",
		"fe":         "falcon-fe",
		"hbs":        "falcon-hbs",
		"judge":      "falcon-judge",
		"nodata":     "falcon-nodata",
		"query":      "falcon-query",
		"sender":     "falcon-sender",
		"task":       "falcon-task",
		"transfer":   "falcon-transfer",
		"gateway":    "falcon-gateway",
	}

	// Modules are deployed in this order
	AllModulesInOrder = []string{
		"graph",
		"hbs",
		"fe",
		"sender",
		"query",
		"judge",
		"transfer",
		"nodata",
		"task",
		"aggregator",
		"agent",
	}
}
