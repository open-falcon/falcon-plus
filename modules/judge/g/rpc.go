package g

import (
	"log"
	"math"
	"sync"
	"time"
)

type SingleConnRpcClient struct {
	sync.Mutex
	rpcClient  *TimeoutRpcClient
	RpcServers []string
	Timeout    time.Duration
}

func (this *SingleConnRpcClient) close() {
	if this.rpcClient != nil {
		this.rpcClient.Close()
		this.rpcClient = nil
	}
}

func (this *SingleConnRpcClient) insureConn() {
	if this.rpcClient != nil {
		return
	}

	var err error
	var retry int = 1

	for {
		if this.rpcClient != nil {
			return
		}

		for _, s := range this.RpcServers {
			this.rpcClient, err = NewTimeoutRpcClient("tcp", s, this.Timeout)
			if err == nil {
				return
			}

			log.Printf("dial %s fail: %s", s, err)
		}

		if retry > 6 {
			retry = 1
		}

		time.Sleep(time.Duration(math.Pow(2.0, float64(retry))) * time.Second)

		retry++
	}
}

func (this *SingleConnRpcClient) Call(method string, args interface{}, reply interface{}) error {

	this.Lock()
	defer this.Unlock()

	this.insureConn()

	err := this.rpcClient.Call(method, args, reply)
	if err != nil {
		this.close()
	}

	return err
}

func (this *SingleConnRpcClient) CallTimeout(method string, args interface{}, reply interface{}, timeout time.Duration) error {

	this.Lock()
	defer this.Unlock()

	this.insureConn()

	err := this.rpcClient.CallTimeout(method, args, reply, timeout)
	if err != nil {
		this.close()
	}

	return err
}
