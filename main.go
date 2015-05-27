package main

import (
	"flag"
	"fmt"
	"github.com/open-falcon/task/collector"
	"github.com/open-falcon/task/g"
	"github.com/open-falcon/task/http"
	"github.com/open-falcon/task/index"
	"github.com/open-falcon/task/monitor"
	"github.com/open-falcon/task/proc"
	"os"
)

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
	// proc
	proc.Start()

	// graph index
	index.Start()
	// collector
	collector.Start()
	// monitor
	monitor.Start()

	// http
	http.Start()

	select {}
}
