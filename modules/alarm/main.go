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
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/open-falcon/falcon-plus/modules/alarm/cron"
	"github.com/open-falcon/falcon-plus/modules/alarm/g"
	"github.com/open-falcon/falcon-plus/modules/alarm/http"
	"github.com/open-falcon/falcon-plus/modules/alarm/model"
)

func main() {
	cfg := flag.String("c", "cfg.json", "configuration file")
	version := flag.Bool("v", false, "show version")
	help := flag.Bool("h", false, "help")
	flag.Parse()

	if *version {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	g.ParseConfig(*cfg)

	g.InitLog(g.Config().LogLevel)
	if g.Config().LogLevel != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}

	g.InitRedisConnPool()
	model.InitDatabase()
	cron.InitSenderWorker()

	go http.Start()

	high_queues := g.Config().Redis.HighQueues
	if len(high_queues) == 0 {
		return
	} else {
		count := len(high_queues)
		for i := 0; i < count; i++ {
			params := make([]interface{}, 2)
			params[0] = high_queues[i]
			params[1] = 0
			go cron.SinglePopEvent(true, params...)
		}

	}

	low_queues := g.Config().Redis.LowQueues
	if len(low_queues) == 0 {
		return
	} else {
		count := len(low_queues)
		for i := 0; i < count; i++ {
			params := make([]interface{}, 2)
			params[0] = low_queues[i]
			params[1] = 0
			go cron.SinglePopEvent(false, params...)
		}

	}

	go cron.CombineSms()
	go cron.CombineMail()
	go cron.CombineIM()
	go cron.ConsumeIM()
	go cron.ConsumeSms()
	go cron.ConsumeMail()
	go cron.CleanExpiredEvent()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		fmt.Println()
		g.RedisClose()
		os.Exit(0)
	}()

	select {}
}
