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
	}

	ModuleBins = map[string]string{
		"agent":      "./bin/falcon-agent",
		"aggregator": "./bin/falcon-aggregator",
		"fe":         "./bin/falcon-fe",
		"graph":      "./bin/falcon-graph",
		"hbs":        "./bin/falcon-hbs",
		"judge":      "./bin/falcon-judge",
		"nodata":     "./bin/falcon-nodata",
		"query":      "./bin/falcon-query",
		"sender":     "./bin/falcon-sender",
		"task":       "./bin/falcon-task",
		"transfer":   "./bin/falcon-transfer",
	}

	ModuleConfs = map[string]string{
		"agent":      "./config/agent.json",
		"aggregator": "./config/aggregator.json",
		"fe":         "./config/api.json",
		"graph":      "./config/graph.json",
		"hbs":        "./config/hbs.json",
		"judge":      "./config/judge.json",
		"nodata":     "./config/nodata.json",
		"query":      "./config/query.json",
		"sender":     "./config/sender.json",
		"task":       "./config/task.json",
		"transfer":   "./config/transfer.json",
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
