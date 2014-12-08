package main

import (
	"flag"
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
}
