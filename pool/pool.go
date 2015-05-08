package pool

import (
	"fmt"
	"github.com/open-falcon/common/model"
	"github.com/toolkits/pool"
	"io"
	"log"
	"net"
	"net/rpc"
	"sync"
	"time"
)

type RpcClient struct {
	cli *rpc.Client
}

func (this RpcClient) Close() (err error) {
	if this.cli != nil {
		err = this.cli.Close()
		this.cli = nil
	}
	return
}

func (this RpcClient) Call(method string, args interface{}, reply interface{}) error {
	return this.cli.Call(method, args, reply)
}

type SafeRpcConnPools struct {
	sync.RWMutex
	M           map[string]*pool.ConnPool
	PingMethod  string
	MaxConns    int
	MaxIdle     int
	ConnTimeout int
	CallTimeout int
}

func CreateSafeRpcConnPools(pingMethod string, maxConns, maxIdle, connTimeout, callTimeout int) *SafeRpcConnPools {
	return &SafeRpcConnPools{
		M:           make(map[string]*pool.ConnPool),
		PingMethod:  pingMethod,
		MaxConns:    maxConns,
		MaxIdle:     maxIdle,
		ConnTimeout: connTimeout,
		CallTimeout: callTimeout,
	}
}

// 同步发送, 完成发送或超时后 才能返回
func (this *SafeRpcConnPools) Call(addr, method string, args interface{}, resp interface{}) error {
	connPool, exists := this.Get(addr)
	if !exists {
		return fmt.Errorf("%s has no connection pool", addr)
	}

	conn, err := connPool.Get()
	if err != nil || conn == nil { //conn可能为空,这也属于异常情况 add by niean
		return fmt.Errorf("%s get connection fail: conn %v, err %v", addr, conn, err)
	}

	rpcClient := conn.(RpcClient)
	callTimeout := time.Duration(this.CallTimeout) * time.Millisecond

	done := make(chan error)
	go func() {
		done <- rpcClient.Call(method, args, resp)
	}()

	select {
	case <-time.After(callTimeout):
		connPool.ForceClose(conn)
		return fmt.Errorf("CallTimeout: %s", addr)
	case err = <-done:
		if err != nil {
			connPool.ForceClose(conn)
		} else {
			connPool.Release(conn)
		}
		return err
	}
}

func (this *SafeRpcConnPools) Exists(address string) bool {
	this.RLock()
	defer this.RUnlock()
	_, exists := this.M[address]
	return exists
}

func (this *SafeRpcConnPools) Get(address string) (*pool.ConnPool, bool) {
	this.RLock()
	defer this.RUnlock()
	p, exists := this.M[address]
	return p, exists
}

func (this *SafeRpcConnPools) Keys() []string {
	this.RLock()
	defer this.RUnlock()
	count := len(this.M)
	keys := make([]string, count)
	i := 0
	for key := range this.M {
		keys[i] = key
		i++
	}
	return keys
}

func (this *SafeRpcConnPools) Put(address string, p *pool.ConnPool) {
	this.Lock()
	defer this.Unlock()
	this.M[address] = p
}

func (this *SafeRpcConnPools) Destroy() {
	this.Lock()
	defer this.Unlock()
	addresses := make([]string, 0, len(this.M))
	for address := range this.M {
		addresses = append(addresses, address)
	}

	for _, address := range addresses {
		this.M[address].Destroy()
		delete(this.M, address)
	}
}

func (this *SafeRpcConnPools) Delete(addr string) {
	p, exists := this.Get(addr)
	if !exists {
		return
	}

	this.Lock()
	delete(this.M, addr)
	this.Unlock()

	p.Destroy()
}

func (this *SafeRpcConnPools) Init(cluster []string) {
	if cluster == nil || len(cluster) == 0 {
		log.Println("no cluster configuration")
		return
	}

	connTimeout := time.Duration(this.ConnTimeout) * time.Millisecond

	for _, address := range cluster {
		if this.Exists(address) {
			continue
		}

		this.M[address] = createOnePool(address, this.PingMethod, connTimeout, this.MaxConns, this.MaxIdle)
	}
}

func createOnePool(address, pingMethod string, connTimeout time.Duration, maxConns, maxIdle int) *pool.ConnPool {
	p := pool.Create(address, maxConns, maxIdle)

	p.New = func() (io.Closer, error) {
		_, err := net.ResolveTCPAddr("tcp", p.Address)
		if err != nil {
			log.Println(p.Address, "format error", err)
			return nil, err
		}

		conn, err := net.DialTimeout("tcp", p.Address, connTimeout)
		if err != nil {
			log.Printf("dial %s fail: %v", p.Address, err)
			return nil, err
		}

		return RpcClient{cli: rpc.NewClient(conn)}, nil
	}

	p.Ping = func(conn io.Closer) error {
		rpcClient := conn.(RpcClient)
		if rpcClient.cli == nil {
			return fmt.Errorf("nil conn")
		}

		resp := &model.SimpleRpcResponse{}
		err := rpcClient.Call(pingMethod, model.NullRpcRequest{}, resp)
		if err != nil {
			log.Println(p.Address, "ping fail", err)
		}
		return err
	}

	p.TestOnBorrow = true

	return p
}
