package cron

import (
	"github.com/open-falcon/falcon-plus/modules/judge/string_matcher"
	"log"
	"time"
)

func CleanStringMatcherHistory() {
	for {
		cleanStringMatcherHistory()
		time.Sleep(time.Hour * 12)
	}
}

func cleanStringMatcherHistory() {
	//aMonthAgo := time.Now().Unix() - 3600*24*7*31
	aMonthAgo := time.Now().Unix() - 300
	err := string_matcher.Consumer.BatchDeleteHistory(aMonthAgo)
	if err != nil {
		log.Println("ERROR: BatchDeleteHistory failed", err)
	}
}
