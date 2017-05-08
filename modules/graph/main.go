package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
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
	log.Println(pid, "register signal notify")
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	for {
		s := <-sigs
		log.Println("recv", s)

		switch s {
		case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			log.Println("gracefull shut down")
			if cfg.Http.Enabled {
				http.Close_chan <- 1
				<-http.Close_done_chan
			}
			log.Println("http stop ok")

			if cfg.Rpc.Enabled {
				api.Close_chan <- 1
				<-api.Close_done_chan
			}
			log.Println("rpc stop ok")

			rrdtool.Out_done_chan <- 1
			rrdtool.FlushAll(true)
			log.Println("rrdtool stop ok")

			log.Println(pid, "exit")
			os.Exit(0)
		}
	}
}

func main() {
	cfg := flag.String("c", "cfg.json", "specify config file")
	version := flag.Bool("v", false, "show version")
	versionGit := flag.Bool("vg", false, "show version and git commit log")
	flag.Parse()

	if *version {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}
	if *versionGit {
		fmt.Println(g.VERSION, g.COMMIT)
		os.Exit(0)
	}

	// global config
	g.ParseConfig(*cfg)

	if g.Config().Debug {
		g.InitLog("debug")
	} else {
		g.InitLog("info")
		gin.SetMode(gin.ReleaseMode)
	}

	// init db
	g.InitDB()
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
