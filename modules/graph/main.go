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
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"

	"github.com/open-falcon/falcon-plus/modules/graph/api"
	"github.com/open-falcon/falcon-plus/modules/graph/cron"
	"github.com/open-falcon/falcon-plus/modules/graph/g"
	"github.com/open-falcon/falcon-plus/modules/graph/http"
	"github.com/open-falcon/falcon-plus/modules/graph/index"
	"github.com/open-falcon/falcon-plus/modules/graph/rrdtool"
)

func start_signal(pid int, cfg *g.GlobalConfig) {
	sigs := make(chan os.Signal, 1)
	log.Info(pid, " register signal notify")
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	for {
		s := <-sigs
		log.Info("recv", s)

		switch s {
		case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			log.Info("graceful shut down")
			if cfg.Http.Enabled {
				http.Close_chan <- 1
				<-http.Close_done_chan
			}
			log.Info("http stop ok")

			if cfg.Rpc.Enabled {
				api.Close_chan <- 1
				<-api.Close_done_chan
			}
			log.Info("rpc stop ok")

			rrdtool.Main_done_chan <- 1
			//flush cache to local file or transmit cache to remote graph
			rrdtool.CommitBeforeQuit()
			log.Info("rrdtool stop ok")

			log.Info("pid:", pid, " exit")
			os.Exit(0)
		}
	}
}

func main() {
	g.BinaryName = BinaryName
	g.Version = Version
	g.GitCommit = GitCommit

	cfg := flag.String("c", "cfg.json", "specify config file")
	version := flag.Bool("v", false, "show version")
	versionGit := flag.Bool("vg", false, "show version and git commit log")
	flag.Parse()

	if *version {
		fmt.Printf("Open-Falcon %s version %s, build %s\n", BinaryName, Version, GitCommit)
		os.Exit(0)
	}
	if *versionGit {
		fmt.Printf("Open-Falcon %s version %s, build %s\n", BinaryName, Version, GitCommit)
		os.Exit(0)
	}

	// global config
	g.ParseConfig(*cfg)

	if g.Config().Debug {
		g.InitLog("debug")
	} else {
		g.InitLog("info")
	}

	// init db
	g.InitDB()
	// rrdtool init
	rrdtool.InitChannel()
	// rrdtool before api for disable loopback connection
	rrdtool.Start()
	// start api
	go api.Start()
	// start indexing
	index.Start()
	// start http server
	go http.Start()
	go cron.CleanCache()

	start_signal(os.Getpid(), g.Config())
}
