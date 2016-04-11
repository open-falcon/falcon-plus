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
	}

	ModuleConfs = map[string]string{
		"agent":      "./agent/config/agent.json",
		"aggregator": "./aggregator/config/aggregator.json",
		"fe":         "./fe/config/api.json",
		"graph":      "./graph/config/graph.json",
		"hbs":        "./hbs/config/hbs.json",
		"judge":      "./judge/config/judge.json",
		"nodata":     "./nodata/config/nodata.json",
		"query":      "./query/config/query.json",
		"sender":     "./sender/config/sender.json",
		"task":       "./task/config/task.json",
		"transfer":   "./transfer/config/transfer.json",
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
