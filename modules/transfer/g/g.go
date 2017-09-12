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

import (
	"log"
	"runtime"
)

// changelog:
// 0.0.1: init project
// 0.0.4: bugfix: set replicas before add node
// 0.0.8: change receiver, mv proc cron to proc pkg, add readme, add gitversion, add config reload, add trace tools
// 0.0.9: fix bugs of conn pool(use transfer's private conn pool, named & minimum)
// 0.0.10: use more efficient proc & sema, rm conn_pool status log
// 0.0.11: fix bug: all graphs' traffic delined when one graph broken down, modify retry interval
// 0.0.14: support sending multi copies to graph node, align ts for judge, add filter
// 0.0.15: support tsdb
// 0.0.16: support config of min step
// 0.0.17: remove migrating, which is implemented in graph

const (
	VERSION      = "0.0.17"
	GAUGE        = "GAUGE"
	COUNTER      = "COUNTER"
	DERIVE       = "DERIVE"
	DEFAULT_STEP = 60
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}
