package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/open-falcon/falcon-plus/modules/gateway/g"
	"github.com/open-falcon/falcon-plus/modules/gateway/http"
	"github.com/open-falcon/falcon-plus/modules/gateway/receiver"
	"github.com/open-falcon/falcon-plus/modules/gateway/sender"
)

func main() {
	cfg := flag.String("c", "cfg.json", "configuration file")
	version := flag.Bool("v", false, "show version")
	flag.Parse()

	if *version {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}

	// global config
	g.ParseConfig(*cfg)

	sender.Start()
	receiver.Start()

	// http
	http.Start()

	select {}
}
