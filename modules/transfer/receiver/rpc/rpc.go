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

package rpc

import (
	"github.com/open-falcon/falcon-plus/modules/transfer/g"
	log "github.com/Sirupsen/logrus"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

func StartRpc() {
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

	server := rpc.NewServer()
	server.Register(new(Transfer))

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("listener.Accept occur error:", err)
			continue
		}
		go server.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}
