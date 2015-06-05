package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/open-falcon/graph/api"
	"github.com/open-falcon/graph/cron"
	"github.com/open-falcon/graph/db"
	"github.com/open-falcon/graph/g"
	"github.com/open-falcon/graph/http"
	"github.com/open-falcon/graph/index"
	"github.com/open-falcon/graph/rrdtool"
)

func start_signal(pid int, conf g.GlobalConfig) {
	sigs := make(chan os.Signal, 1)
	log.Println(pid, "register signal notify")
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	for {
		s := <-sigs
		log.Println("recv", s)

		switch s {
		case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			log.Println("gracefull shut down")
			http.Close_chan <- 1
			api.Close_chan <- 1
			cron.Out_done_chan <- 1
			<-http.Close_done_chan
			<-api.Close_done_chan
			log.Println(pid, "remove ", conf.Pid)
			os.Remove(conf.Pid)
			rrdtool.FlushAll()
			log.Println(pid, "flush done, exiting")
			os.Exit(0)
		}
	}
}

func main() {

	cfg := flag.String("c", "cfg.json", "configuration file")
	version := flag.Bool("v", false, "show version")
	versionGit := flag.Bool("vg", false, "show version")
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

	// init db
	db.Init()
	// start rrdtool
	rrdtool.Start()

	go api.Start()

	// 刷硬盘
	go cron.SyncDisk()

	// 索引更新2.0
	index.Start()

	// http
	go http.Start()

	start_signal(os.Getpid(), *g.Config())

}
