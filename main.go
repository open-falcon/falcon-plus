package main

import (
	"flag"
	"fmt"
	"github.com/open-falcon/query/g"
	"github.com/open-falcon/query/graph"
	"github.com/open-falcon/query/http"
	"github.com/toolkits/logger"
	"os"
	"os/signal"
	"syscall"
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

	// config
	g.ParseConfig(*cfg)

	var err error
	// graph section
	err = graph.InitBackends()
	if err != nil {
		logger.Error("load graph backends fail: %v", err)
	}

	go graph.ReloadBackends()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		fmt.Println("END")
		graph.DestroyConnPools()
		os.Exit(0)
	}()

	http.Start()
}
