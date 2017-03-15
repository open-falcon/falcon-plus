package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/open-falcon/falcon-plus/modules/ctrl/falcon"
	"github.com/open-falcon/falcon-plus/modules/transfer/g"
	"github.com/open-falcon/falcon-plus/modules/transfer/http"
	"github.com/open-falcon/falcon-plus/modules/transfer/proc"
	"github.com/open-falcon/falcon-plus/modules/transfer/receiver"
	"github.com/open-falcon/falcon-plus/modules/transfer/sender"
)

func leaseStart() {
	c := falcon.NewEtcdCli2(g.Config().Lease)
	c.Prestart()
	c.Start()
}

func main() {
	cfg := flag.String("c", "cfg.json", "configuration file")
	version := flag.Bool("version", false, "show version")
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
	// lease key to etcd server (v3)
	leaseStart()
	// proc
	proc.Start()

	sender.Start()
	receiver.Start()

	// http
	http.Start()

	select {}
}
