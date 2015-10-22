package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/open-falcon/nodata/collector"
	"github.com/open-falcon/nodata/config"
	"github.com/open-falcon/nodata/g"
	"github.com/open-falcon/nodata/http"
	"github.com/open-falcon/nodata/judge"
	"github.com/open-falcon/nodata/sender"
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
	g.StartProc()

	// config
	config.Start()
	// collector
	collector.Start()
	// judge
	judge.Start()
	// sender
	sender.Start()

	// http
	http.Start()

	select {}
}
