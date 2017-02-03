package cron

import (
	"github.com/open-falcon/falcon-plus/modules/judge/store"
	"time"
)

func CleanStale() {
	for {
		time.Sleep(time.Hour * 5)
		cleanStale()
	}
}

func cleanStale() {
	before := time.Now().Unix() - 3600*24*7

	arr := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f"}
	for i := 0; i < 16; i++ {
		for j := 0; j < 16; j++ {
			store.HistoryBigMap[arr[i]+arr[j]].CleanStale(before)
		}
	}
}
