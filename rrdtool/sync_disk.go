package rrdtool

import (
	"io/ioutil"
	"log"
	"time"

	"github.com/open-falcon/graph/g"
	"github.com/open-falcon/graph/store"
)

const (
	_ = iota
	IO_TASK_M_READ
	IO_TASK_M_FLUSH
	IO_TASK_M_FETCH
)

type io_task_t struct {
	method int
	args   interface{}
	done   chan error
}

var (
	Out_done_chan chan int
	io_task_chan  chan *io_task_t
)

func init() {
	Out_done_chan = make(chan int, 1)
	io_task_chan = make(chan *io_task_t, 16)
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

func ioWorker() {
	var err error
	for {
		select {
		case task := <-io_task_chan:
			if task.method == IO_TASK_M_READ {
				if args, ok := task.args.(*readfile_t); ok {
					args.data, err = ioutil.ReadFile(args.filename)
					task.done <- err
				}
			} else if task.method == IO_TASK_M_FLUSH {
				if args, ok := task.args.(*flushfile_t); ok {
					task.done <- flushrrd(args.filename, args.items)
				}
			} else if task.method == IO_TASK_M_FETCH {
				if args, ok := task.args.(*fetch_t); ok {
					args.data, err = fetch(args.filename, args.cf, args.start, args.end, args.step)
					task.done <- err
				}
			}
		}
	}
}
