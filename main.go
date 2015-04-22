package main

import (
	"flag"
	"fmt"
	"github.com/open-falcon/task/db"
	"github.com/open-falcon/task/g"
	"github.com/open-falcon/task/http"
	"github.com/open-falcon/task/index"
	"github.com/open-falcon/task/monitor"
	"github.com/open-falcon/task/proc"
	"os"
	"os/signal"
	"syscall"
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

	// global config
	g.ParseConfig(*cfg)
	// proc
	proc.Init()
	// db
	db.Init()

	// graph index
	index.Start()
	// monitor
	monitor.Start()

	// http
	// ENSURE starting httpServer to be the final step
	http.Start()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		fmt.Println()
		os.Exit(0)
	}()

	select {}
}
