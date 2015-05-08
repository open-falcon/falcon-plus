package cron

import (
	"log"
)

func Start() {
	// proc
	go StartProcCron()
	log.Println("cron.Start, ok")
}
