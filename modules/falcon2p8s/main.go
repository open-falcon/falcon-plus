package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/open-falcon/falcon-plus/modules/falcon2p8s/g"
	"github.com/open-falcon/falcon-plus/modules/falcon2p8s/http"
	"github.com/open-falcon/falcon-plus/modules/falcon2p8s/rpc"
	log "github.com/sirupsen/logrus"
)

func startSignal(pid int) {
	sigs := make(chan os.Signal, 1)
	log.Println(pid, "register signal notify")
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	for {
		s := <-sigs
		log.Println("recv", s)
		switch s {
		case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			for {
				if g.P8sItemQueue.Empty() {
					log.Println("graceful shut down")
					log.Println(pid, "exit")
					os.Exit(0)
				}
			}
		}
	}
}

func main() {
	g.BinaryName = BinaryName
	g.Version = Version
	g.GitCommit = GitCommit

	cfg := flag.String("c", "cfg.json", "configuration file")
	version := flag.Bool("v", false, "show version")
	help := flag.Bool("h", false, "help")
	flag.Parse()

	if *version {
		fmt.Printf("Open-Falcon %s version %s, build %s\n", BinaryName, Version, GitCommit)
		os.Exit(0)
	}

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	g.ParseConfig(*cfg)

	g.InitLog(g.Config().LogLevel)

	go rpc.Start()
	go http.Start()
	startSignal(os.Getpid())
}
