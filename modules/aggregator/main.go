package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/open-falcon/falcon-plus/common/sdk/graph"
	"github.com/open-falcon/falcon-plus/common/sdk/portal"
	"github.com/open-falcon/falcon-plus/common/sdk/sender"
	"github.com/open-falcon/falcon-plus/modules/aggregator/cron"
	"github.com/open-falcon/falcon-plus/modules/aggregator/db"
	"github.com/open-falcon/falcon-plus/modules/aggregator/g"
	"github.com/open-falcon/falcon-plus/modules/aggregator/http"
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
	db.Init()

	go http.Start()
	go cron.UpdateItems()

	// sdk configuration
	graph.GraphLastUrl = g.Config().Api.GraphLast
	sender.Debug = g.Config().Debug
	sender.PostPushUrl = g.Config().Api.Push
	portal.HostnamesUrl = g.Config().Api.Hostnames

	sender.StartSender()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		fmt.Println()
		os.Exit(0)
	}()

	select {}
}
