package rpc

import (
	"github.com/open-falcon/judge/g"
	"log"
	"net"
	"net/rpc"
)

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

	rpc.Register(new(Judge))

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("listener.Accept occur error: %s", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}
