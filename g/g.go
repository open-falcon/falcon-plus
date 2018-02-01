// Copyright 2017 Xiaomi, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package g

import "path/filepath"

var Modules map[string]bool
var BinOf map[string]string
var cfgOf map[string]string
var ModuleApps map[string]string
var logpathOf map[string]string
var PidOf map[string]string
var AllModulesInOrder []string

func init() {
	Modules = map[string]bool{
		"agent":      true,
		"aggregator": true,
		"graph":      true,
		"hbs":        true,
		"judge":      true,
		"nodata":     true,
		"transfer":   true,
		"gateway":    true,
		"api":        true,
		"alarm":      true,
	}

	BinOf = map[string]string{
		"agent":      "./bin/falcon-agent",
		"aggregator": "./bin/falcon-aggregator",
		"graph":      "./bin/falcon-graph",
		"hbs":        "./bin/falcon-hbs",
		"judge":      "./bin/falcon-judge",
		"nodata":     "./bin/falcon-nodata",
		"transfer":   "./bin/falcon-transfer",
		"gateway":    "./bin/falcon-gateway",
		"api":        "./bin/falcon-api",
		"alarm":      "./bin/falcon-alarm",
	}

	cfgOf = map[string]string{
		"agent":      "./config/agent-cfg.json",
		"aggregator": "./config/aggregator-cfg.json",
		"graph":      "./config/graph-cfg.json",
		"hbs":        "./config/hbs-cfg.json",
		"judge":      "./config/judge-cfg.json",
		"nodata":     "./config/nodata-cfg.json",
		"transfer":   "./config/transfer-cfg.json",
		"gateway":    "./config/gateway-cfg.json",
		"api":        "./config/api-cfg.json",
		"alarm":      "./config/alarm-cfg.json",
	}

	ModuleApps = map[string]string{
		"agent":      "falcon-agent",
		"aggregator": "falcon-aggregator",
		"graph":      "falcon-graph",
		"hbs":        "falcon-hbs",
		"judge":      "falcon-judge",
		"nodata":     "falcon-nodata",
		"transfer":   "falcon-transfer",
		"gateway":    "falcon-gateway",
		"api":        "falcon-api",
		"alarm":      "falcon-alarm",
	}

	logpathOf = map[string]string{
		"agent":      "./logs/agent.log",
		"aggregator": "./logs/aggregator.log",
		"graph":      "./logs/graph.log",
		"hbs":        "./logs/hbs.log",
		"judge":      "./logs/judge.log",
		"nodata":     "./logs/nodata.log",
		"transfer":   "./logs/transfer.log",
		"gateway":    "./logs/gateway.log",
		"api":        "./logs/api.log",
		"alarm":      "./logs/alarm.log",
	}

	PidOf = map[string]string{
		"agent":      "<NOT SET>",
		"aggregator": "<NOT SET>",
		"graph":      "<NOT SET>",
		"hbs":        "<NOT SET>",
		"judge":      "<NOT SET>",
		"nodata":     "<NOT SET>",
		"transfer":   "<NOT SET>",
		"gateway":    "<NOT SET>",
		"api":        "<NOT SET>",
		"alarm":      "<NOT SET>",
	}

	// Modules are deployed in this order
	AllModulesInOrder = []string{
		"graph",
		"hbs",
		"judge",
		"transfer",
		"nodata",
		"aggregator",
		"agent",
		"gateway",
		"api",
		"alarm",
	}
}

func Bin(name string) string {
	p, _ := filepath.Abs(BinOf[name])
	return p
}

func Cfg(name string) string {
	p, _ := filepath.Abs(cfgOf[name])
	return p
}

func LogPath(name string) string {
	p, _ := filepath.Abs(logpathOf[name])
	return p
}

func LogDir(name string) string {
	d, _ := filepath.Abs(filepath.Dir(logpathOf[name]))
	return d
}
