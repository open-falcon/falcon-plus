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

package api

import (
	"container/list"
	"log"
	"net"
	"net/rpc"
	"sync"
	"time"

	"github.com/open-falcon/falcon-plus/modules/graph/g"
)

type conn_list struct {
	sync.RWMutex
	list *list.List
}

func (l *conn_list) insert(c net.Conn) *list.Element {
	l.Lock()
	defer l.Unlock()
	return l.list.PushBack(c)
}
func (l *conn_list) remove(e *list.Element) net.Conn {
	l.Lock()
	defer l.Unlock()
	return l.list.Remove(e).(net.Conn)
}

var Close_chan, Close_done_chan chan int
var connects conn_list

func init() {
	Close_chan = make(chan int, 1)
	Close_done_chan = make(chan int, 1)
	connects = conn_list{list: list.New()}
}

func Start() {
	if !g.Config().Rpc.Enabled {
		log.Println("rpc.Start warning, not enabled")
		return
	}
	addr := g.Config().Rpc.Listen
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		log.Fatalf("rpc.Start error, net.ResolveTCPAddr failed, %s", err)
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatalf("rpc.Start error, listen %s failed, %s", addr, err)
	} else {
		log.Println("rpc.Start ok, listening on", addr)
	}

	rpc.Register(new(Graph))

	go func() {
		var tempDelay time.Duration // how long to sleep on accept failure
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				time.Sleep(tempDelay)
				continue
			}

			conn.SetKeepAlive(true)
			tempDelay = 0
			go func() {
				e := connects.insert(conn)
				defer connects.remove(e)
				rpc.ServeConn(conn)
			}()
		}
	}()

	select {
	case <-Close_chan:
		log.Println("rpc, recv sigout and exiting...")
		listener.Close()
		Close_done_chan <- 1

		connects.Lock()
		for e := connects.list.Front(); e != nil; e = e.Next() {
			e.Value.(net.Conn).Close()
		}
		connects.Unlock()

		return
	}

}
