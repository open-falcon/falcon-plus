package main

import (
	"flag"
	"github.com/open-falcon/agent/g"
	"log"
)

func main() {
	var cfg string
	flag.StringVar(&cfg, "c", "", "configuration file")
	flag.Parse()

	if cfg == "" {
		log.Fatalln("use -c to specify configuration file")
	}

	g.ParseConfig(cfg)

	log.Println(g.GetConfig().Heartbeat.Addr)
}
