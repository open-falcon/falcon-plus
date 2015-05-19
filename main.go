package main

import (
	"flag"
	"fmt"
	"github.com/open-falcon/transfer/g"
	"github.com/open-falcon/transfer/http"
	"github.com/open-falcon/transfer/proc"
	"github.com/open-falcon/transfer/receiver"
	"github.com/open-falcon/transfer/sender"
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

	sender.Start()
	receiver.Start()

	// http
	http.Start()

	select {}
}
