package cron

import (
	"log"
	"time"

	"github.com/open-falcon/graph/g"
	"github.com/open-falcon/graph/rrdtool"
	"github.com/open-falcon/graph/store"
)

var (
	Out_done_chan chan int
	Counter       uint64
)

func init() {
	Counter = 0
	Out_done_chan = make(chan int, 1)
}

func SyncDisk() {
	time.Sleep(time.Second * 300)
	ticker := time.NewTicker(time.Millisecond * g.FLUSH_DISK_STEP).C
	var idx int = 0

	for {
		select {
		case <-ticker:
			idx = idx % store.GraphItems.Size
			rrdtool.FlushRRD(idx)
			idx += 1
		case <-Out_done_chan:
			log.Println("cron recv sigout and exit...")
			return
		}
	}
}
