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
	"github.com/toolkits/logger"
)

var pid int

var conf g.GlobalConfig

func start_signal() {
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
	pid = os.Getpid()

	cfg := flag.String("c", "cfg.json", "configuration file")
	version := flag.Bool("v", false, "show version")
	versionGit := flag.Bool("vg", false, "show version")
	flag.Parse()
	g.ParseConfig(*cfg)

	conf = *g.Config()

	if *version {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}

	if *versionGit {
		fmt.Println(g.VERSION, g.COMMIT)
		os.Exit(0)
	}

	// 只在启动的时候初始化一次，而ParseConfig可以被多次调用
	logger.SetLevelWithDefault(g.Config().Log, "info")
	db.Init()

	go api.Start()

	// 刷硬盘
	go cron.SyncDisk()

	// 索引更新2.0
	index.Start()

	// http
	go http.Start()

	start_signal()

}
