package main

import (
	"flag"
	"fmt"
	"github.com/open-falcon/agent/funcs"
	"github.com/open-falcon/agent/g"
	"os"
)

func main() {

	cfg := flag.String("c", "cfg.json", "configuration file")
	version := flag.Bool("v", false, "show version")
	print := flag.Bool("p", false, "print all metrics")

	flag.Parse()

	if *version {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}

	if *print {
		funcs.PrintAll()
		os.Exit(0)
	}

	g.ParseConfig(*cfg)
	g.InitVars()

}
