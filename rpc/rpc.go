package rpc

import (
	"github.com/open-falcon/hbs/g"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type Hbs int
type Agent int

func Start() {
	addr := g.Config().Listen

	server := rpc.NewServer()
	// server.Register(new(filter.Filter))
	server.Register(new(Agent))
	server.Register(new(Hbs))

	l, e := net.Listen("tcp", addr)
	if e != nil {
		log.Fatalln("listen error:", e)
	} else {
		log.Println("listening", addr)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("listener accept fail:", err)
			continue
		}
		go server.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}
