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

	g.ParseConfig(cfg)

	log.Println(g.Config().Http.Port)
}
