package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/open-falcon/falcon-plus/modules/alarm/cron"
	"github.com/open-falcon/falcon-plus/modules/alarm/g"
	"github.com/open-falcon/falcon-plus/modules/alarm/http"
	"github.com/open-falcon/falcon-plus/modules/alarm/model"
)

func main() {
	cfg := flag.String("c", "cfg.json", "configuration file")
	version := flag.Bool("v", false, "show version")
	help := flag.Bool("h", false, "help")
	flag.Parse()

	if *version {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	g.ParseConfig(*cfg)

	g.InitLog(g.Config().LogLevel)
	if g.Config().LogLevel != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}

	g.InitRedisConnPool()
	model.InitDatabase()
	cron.InitSenderWorker()

	go http.Start()
	go cron.ReadHighEvent()
	go cron.ReadLowEvent()
	go cron.CombineChat()
	go cron.CombineSms()
	go cron.CombineMail()
	go cron.ConsumeChat()
	go cron.ConsumeSms()
	go cron.ConsumeMail()
	go cron.CleanExpiredEvent()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		fmt.Println()
		g.RedisConnPool.Close()
		os.Exit(0)
	}()

	select {}
}
