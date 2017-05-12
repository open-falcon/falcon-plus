package backend_pool

import (
	"fmt"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"sync"
	"time"

	connp "github.com/toolkits/conn_pool"
	rpcpool "github.com/toolkits/conn_pool/rpc_conn_pool"
)

// ConnPools Manager
type SafeRpcConnPools struct {
	sync.RWMutex
	M           map[string]*connp.ConnPool
	MaxConns    int
	MaxIdle     int
	ConnTimeout int
	CallTimeout int
}

func CreateSafeRpcConnPools(maxConns, maxIdle, connTimeout, callTimeout int, cluster []string) *SafeRpcConnPools {
	cp := &SafeRpcConnPools{M: make(map[string]*connp.ConnPool), MaxConns: maxConns, MaxIdle: maxIdle,
		ConnTimeout: connTimeout, CallTimeout: callTimeout}

	ct := time.Duration(cp.ConnTimeout) * time.Millisecond
	for _, address := range cluster {
		if _, exist := cp.M[address]; exist {
			continue
		}
		cp.M[address] = createOneRpcPool(address, address, ct, maxConns, maxIdle)
	}

	return cp
}

func CreateSafeJsonrpcConnPools(maxConns, maxIdle, connTimeout, callTimeout int, cluster []string) *SafeRpcConnPools {
	cp := &SafeRpcConnPools{M: make(map[string]*connp.ConnPool), MaxConns: maxConns, MaxIdle: maxIdle,
		ConnTimeout: connTimeout, CallTimeout: callTimeout}

	ct := time.Duration(cp.ConnTimeout) * time.Millisecond
	for _, address := range cluster {
		if _, exist := cp.M[address]; exist {
			continue
		}
		cp.M[address] = createOneJsonrpcPool(address, address, ct, maxConns, maxIdle)
	}

	return cp
}

// 同步发送, 完成发送或超时后 才能返回
func (this *SafeRpcConnPools) Call(addr, method string, args interface{}, resp interface{}) error {
	connPool, exists := this.Get(addr)
	if !exists {
		return fmt.Errorf("%s has no connection pool", addr)
	}

	conn, err := connPool.Fetch()
	if err != nil {
		return fmt.Errorf("%s get connection fail: conn %v, err %v. proc: %s", addr, conn, err, connPool.Proc())
	}

	rpcClient := conn.(*rpcpool.RpcClient)
	callTimeout := time.Duration(this.CallTimeout) * time.Millisecond

	done := make(chan error, 1)
	go func() {
		done <- rpcClient.Call(method, args, resp)
	}()

	select {
	case <-time.After(callTimeout):
		connPool.ForceClose(conn)
		return fmt.Errorf("%s, call timeout", addr)
	case err = <-done:
		if err != nil {
			connPool.ForceClose(conn)
			err = fmt.Errorf("%s, call failed, err %v. proc: %s", addr, err, connPool.Proc())
		} else {
			connPool.Release(conn)
		}
		return err
	}
}

func (this *SafeRpcConnPools) Get(address string) (*connp.ConnPool, bool) {
	this.RLock()
	defer this.RUnlock()
	p, exists := this.M[address]
	return p, exists
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

func (this *SafeRpcConnPools) Proc() []string {
	procs := []string{}
	for _, cp := range this.M {
		procs = append(procs, cp.Proc())
	}
	return procs
}

func createOneRpcPool(name string, address string, connTimeout time.Duration, maxConns int, maxIdle int) *connp.ConnPool {
	p := connp.NewConnPool(name, address, int32(maxConns), int32(maxIdle))
	p.New = func(connName string) (connp.NConn, error) {
		_, err := net.ResolveTCPAddr("tcp", p.Address)
		if err != nil {
			return nil, err
		}

		conn, err := net.DialTimeout("tcp", p.Address, connTimeout)
		if err != nil {
			return nil, err
		}

		return rpcpool.NewRpcClient(rpc.NewClient(conn), connName), nil
	}

	return p
}

func createOneJsonrpcPool(name string, address string, connTimeout time.Duration, maxConns int, maxIdle int) *connp.ConnPool {
	p := connp.NewConnPool(name, address, int32(maxConns), int32(maxIdle))
	p.New = func(connName string) (connp.NConn, error) {
		_, err := net.ResolveTCPAddr("tcp", p.Address)
		if err != nil {
			return nil, err
		}

		conn, err := net.DialTimeout("tcp", p.Address, connTimeout)
		if err != nil {
			return nil, err
		}

		return rpcpool.NewRpcClientWithCodec(jsonrpc.NewClientCodec(conn), connName), nil
	}

	return p
}
