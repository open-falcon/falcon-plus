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

package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/open-falcon/falcon-plus/modules/judge/cron"
	"github.com/open-falcon/falcon-plus/modules/judge/g"
	"github.com/open-falcon/falcon-plus/modules/judge/http"
	"github.com/open-falcon/falcon-plus/modules/judge/model"
	"github.com/open-falcon/falcon-plus/modules/judge/rpc"
	"github.com/open-falcon/falcon-plus/modules/judge/store"
	"github.com/open-falcon/falcon-plus/modules/judge/string_matcher"
)

func main() {
	cfg := flag.String("c", "cfg.json", "configuration file")
	version := flag.Bool("v", false, "show version")
	flag.Parse()

	if *version {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}

	g.ParseConfig(*cfg)
	_cfg := g.Config()

	if _cfg.Debug {
		log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	}

	g.InitRedisConnPool()
	g.InitHbsClient()

	// string matcher history persistence
	if _cfg.StringMatcher != nil && _cfg.StringMatcher.Enabled {
		string_matcher.InitStringMatcher()
		batch := _cfg.StringMatcher.Batch
		retry := _cfg.StringMatcher.MaxRetry
		go string_matcher.Consumer.Start(batch, retry)
	}

	store.InitHistoryBigMap()
	model.InitEHistoryBigMap()

	go http.Start()
	go rpc.Start()

	go cron.SyncStrategies()
	go cron.SyncEExps()
	go cron.CleanStale()

	if _cfg.StringMatcher != nil && _cfg.StringMatcher.Enabled {
		go cron.CleanStringMatcherHistory()
	}

	select {}
}
