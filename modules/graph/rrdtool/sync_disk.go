// Copyright 2017 Xiaomi, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rrdtool

import (
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/open-falcon/falcon-plus/modules/graph/g"
	"github.com/open-falcon/falcon-plus/modules/graph/store"
	log "github.com/sirupsen/logrus"
	"github.com/toolkits/file"
)

const (
	_ = iota
	IO_TASK_M_READ
	IO_TASK_M_WRITE
	IO_TASK_M_FLUSH
	IO_TASK_M_FETCH
)

type io_task_t struct {
	method int
	args   interface{}
	done   chan error
}

var (
	Main_done_chan chan int
	io_task_chans  []chan *io_task_t
)

func InitChannel() {
	Main_done_chan = make(chan int, 1)
	ioWorkerNum := g.Config().IOWorkerNum
	io_task_chans = make([]chan *io_task_t, ioWorkerNum)
	for i := 0; i < ioWorkerNum; i++ {
		//the io task queue length is 16
		io_task_chans[i] = make(chan *io_task_t, 16)
	}
}

func syncDisk() {
	ticker := time.NewTicker(time.Millisecond * g.FLUSH_DISK_STEP)
	defer ticker.Stop()
	var idx int = 0
	n := store.GraphItems.Size / 10

	for {
		select {
		case <-ticker.C:
			idx = idx % store.GraphItems.Size
			commitByIdx(idx)
			if idx%n == 0 {
				log.Debugf("flush rrd hash idx:%03d size:%03d disk:%08d net:%08d",
					idx, store.GraphItems.Size, disk_counter, net_counter)
			}
			idx += 1
		case <-Main_done_chan:
			log.Info("syncDisk cron recv sigout and exit...")
			return
		}
	}
}

// WriteFile writes data to a file named by filename.
// file must not exist
func writeFile(filename string, data []byte, perm os.FileMode) error {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_EXCL, perm)
	if err != nil {
		return err
	}
	n, err := f.Write(data)
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	}
	if err1 := f.Close(); err == nil {
		err = err1
	}
	return err
}

func ioWorker() {
	ioWorkerNum := g.Config().IOWorkerNum
	for i := 0; i < ioWorkerNum; i++ {
		go func(i int) {
			var err error
			for {
				select {
				case task := <-io_task_chans[i]:
					if task.method == IO_TASK_M_READ {
						if args, ok := task.args.(*readfile_t); ok {
							args.data, err = ioutil.ReadFile(args.filename)
							task.done <- err
						}
					} else if task.method == IO_TASK_M_WRITE {
						//filename must not exist
						if args, ok := task.args.(*g.File); ok {
							baseDir := file.Dir(args.Filename)
							if err = file.InsureDir(baseDir); err != nil {
								task.done <- err
							}
							task.done <- writeFile(args.Filename, args.Body, 0644)
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
		}(i)
	}
}
