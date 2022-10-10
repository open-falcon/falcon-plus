package rpc

import (
	"net"
	"net/rpc"

	"github.com/open-falcon/falcon-plus/modules/falcon2p8s/g"
	log "github.com/sirupsen/logrus"
)

func Start() {
	addr := g.Config().Rpc.Listen
	if addr == "" {
		addr = "0.0.0.0:8080"
	}
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

	rpc.Register(new(P8sRelay))

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("listener.Accept occur error: %s", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}
