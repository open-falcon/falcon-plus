package rrdtool

import (
	"log"
	"time"

	"github.com/open-falcon/graph/g"
	"github.com/open-falcon/graph/store"
)

var (
	Out_done_chan chan int
)

func init() {
	Out_done_chan = make(chan int, 1)
}

func syncDisk() {
	time.Sleep(time.Second * 300)
	ticker := time.NewTicker(time.Millisecond * g.FLUSH_DISK_STEP).C
	var idx int = 0

	for {
		select {
		case <-ticker:
			idx = idx % store.GraphItems.Size
			FlushRRD(idx)
			idx += 1
		case <-Out_done_chan:
			log.Println("cron recv sigout and exit...")
			return
		}
	}
}
