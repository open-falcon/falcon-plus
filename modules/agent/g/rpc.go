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

package g

import (
	"errors"
	"github.com/toolkits/net"
	"log"
	"math"
	"net/rpc"
	"sync"
	"time"
)

type SingleConnRpcClient struct {
	sync.Mutex
	rpcClient  *rpc.Client
	RpcServers []string
	Timeout    time.Duration
}

func (this *SingleConnRpcClient) close() {
	if this.rpcClient != nil {
		this.rpcClient.Close()
		this.rpcClient = nil
	}
}

func (this *SingleConnRpcClient) serverConn() error {
	if this.rpcClient != nil {
		return nil
	}

	var err error
	var retry int

	for _, addr := range this.RpcServers {
		retry = 1

	RETRY:
		this.rpcClient, err = net.JsonRpcClient("tcp", addr, this.Timeout)
		if err != nil {
			log.Println("net.JsonRpcClient failed", err)
			if retry > 3 {
				continue
			}

			time.Sleep(time.Duration(math.Pow(2.0, float64(retry))) * time.Second)
			retry++
			goto RETRY
		}
		log.Println("connected RPC server", addr)

		return nil
	}

	return errors.New("connect to RPC servers failed")
}

func (this *SingleConnRpcClient) Call(method string, args interface{}, reply interface{}) error {

	this.Lock()
	defer this.Unlock()

	err := this.serverConn()
	if err != nil {
		return err
	}

	timeout := time.Duration(10 * time.Second)
	done := make(chan error, 1)

	go func() {
		err := this.rpcClient.Call(method, args, reply)
		done <- err
	}()

	select {
	case <-time.After(timeout):
		log.Printf("[WARN] rpc call timeout %v => %v", this.rpcClient, this.RpcServers)
		this.close()
		return errors.New("rpc call timeout")
	case err := <-done:
		if err != nil {
			this.close()
			return err
		}
	}

	return nil
}
