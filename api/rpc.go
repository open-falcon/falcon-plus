package api

import (
	"log"
	"net"
	"net/rpc"
	"time"

	"github.com/open-falcon/graph/g"
)

var Close_chan, Close_done_chan chan int

func init() {
	Close_chan = make(chan int, 1)
	Close_done_chan = make(chan int, 1)
}

func Start() {
	if !g.Config().Rpc.Enabled {
		return
	}
	addr := g.Config().Rpc.Listen
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		log.Fatalf("net.ResolveTCPAddr fail: %s", err)
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatalf("listen %s fail: %s", addr, err)
	} else {
		log.Println("rpc listening", addr)
	}

	rpc.Register(new(Graph))

	go func() {
		var tempDelay time.Duration // how long to sleep on accept failure
		for {
			conn, err := listener.Accept()
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
			tempDelay = 0
			go rpc.ServeConn(conn)
		}
	}()

	select {
	case <-Close_chan:
		log.Println("api recv sigout and exit...")
		listener.Close()
		Close_done_chan <- 1
		return
	}

}
