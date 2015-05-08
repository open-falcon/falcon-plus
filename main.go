package main

import (
	"flag"
	"fmt"
	"github.com/open-falcon/transfer/cron"
	"github.com/open-falcon/transfer/g"
	"github.com/open-falcon/transfer/http"
	"github.com/open-falcon/transfer/proc"
	"github.com/open-falcon/transfer/rpc"
	"github.com/open-falcon/transfer/sender"
	"github.com/open-falcon/transfer/socket"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := flag.String("c", "cfg.json", "configuration file")
	version := flag.Bool("v", false, "show version")
	versionGit := flag.Bool("vg", false, "show version")
	help := flag.Bool("h", false, "help")
	flag.Parse()

	if *version {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}

	if *versionGit {
		fmt.Println(g.VERSION, g.COMMIT)
		os.Exit(0)
	}

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	// global config
	g.ParseConfig(*cfg)
	// proc
	proc.Init()

	sender.Start()

	go rpc.Start()
	go socket.Start()
	cron.Start()

	// http
	// ENSURE that starting httpServer is the FINAL STEP
	http.Start()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		fmt.Println()
		sender.DestroyConnPools()
		os.Exit(0)
	}()

	select {}
}
