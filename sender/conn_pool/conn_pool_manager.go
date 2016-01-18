package conn_pool

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"sync"
	"time"
)

// RpcCient, 要实现io.Closer接口
type RpcClient struct {
	cli  *rpc.Client
	name string
}

func (this RpcClient) Name() string {
	return this.name
}

func (this RpcClient) Closed() bool {
	return this.cli == nil
}

func (this RpcClient) Close() error {
	if this.cli != nil {
		err := this.cli.Close()
		this.cli = nil
		return err
	}
	return nil
}

func (this RpcClient) Call(method string, args interface{}, reply interface{}) error {
	return this.cli.Call(method, args, reply)
}

// ConnPools Manager
type SafeRpcConnPools struct {
	sync.RWMutex
	M           map[string]*ConnPool
	MaxConns    int
	MaxIdle     int
	ConnTimeout int
	CallTimeout int
}

func CreateSafeRpcConnPools(maxConns, maxIdle, connTimeout, callTimeout int, cluster []string) *SafeRpcConnPools {
	cp := &SafeRpcConnPools{M: make(map[string]*ConnPool), MaxConns: maxConns, MaxIdle: maxIdle,
		ConnTimeout: connTimeout, CallTimeout: callTimeout}

	ct := time.Duration(cp.ConnTimeout) * time.Millisecond
	for _, address := range cluster {
		if _, exist := cp.M[address]; exist {
			continue
		}
		cp.M[address] = createOnePool(address, address, ct, maxConns, maxIdle)
	}

	return cp
}

func (this *SafeRpcConnPools) Proc() []string {
	procs := []string{}
	for _, cp := range this.M {
		procs = append(procs, cp.Proc())
	}
	return procs
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

	rpcClient := conn.(RpcClient)
	callTimeout := time.Duration(this.CallTimeout) * time.Millisecond

	done := make(chan error)
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

func (this *SafeRpcConnPools) Get(address string) (*ConnPool, bool) {
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

func createOnePool(name string, address string, connTimeout time.Duration, maxConns int, maxIdle int) *ConnPool {
	p := NewConnPool(name, address, maxConns, maxIdle)
	p.New = func(connName string) (NConn, error) {
		_, err := net.ResolveTCPAddr("tcp", p.Address)
		if err != nil {
			//log.Println(p.Address, "format error", err)
			return nil, err
		}

		conn, err := net.DialTimeout("tcp", p.Address, connTimeout)
		if err != nil {
			//log.Printf("new conn fail, addr %s, err %v", p.Address, err)
			return nil, err
		}

		return RpcClient{cli: rpc.NewClient(conn), name: connName}, nil
	}

	return p
}

type SafeTcpConnPools struct {
	sync.RWMutex
	M           chan net.Conn
	MaxSize     int
	ConnTimeout int
	CallTimeout int
	MaxRetry    int
	Address     string
}

func (this *SafeTcpConnPools) Send(data []byte) (err error) {
	this.Lock()
	defer this.Unlock()

	conn := <-this.M
	defer func() {
		this.M <- conn
	}()

	if conn == nil {
		conn, err = createOneTcpConn(this.Address, this.ConnTimeout, this.MaxRetry)
		if err != nil {
			return
		}
	}

	done := make(chan error)
	go func() {
		_, err = conn.Write(data)
		done <- err
	}()

	select {
	case <-time.After(time.Duration(this.CallTimeout) * time.Millisecond):
		conn.Close()
		conn = nil
		err = fmt.Errorf("%s, call timeout", this.Address)
	case err = <-done:
		if err != nil {
			conn.Close()
			conn = nil
			err = fmt.Errorf("%s send data fail: %v", this.Address, err)
		}
	}

	return
}

func CreateSafeTcpConnPools(maxRetry, maxSize, connTimeout, callTimeout int, address string) *SafeTcpConnPools {
	cp := &SafeTcpConnPools{M: make(chan net.Conn, maxSize), MaxRetry: maxRetry, MaxSize: maxSize,
		ConnTimeout: connTimeout, CallTimeout: callTimeout, Address: address}

	for i := 0; i < maxSize; i++ {
		conn, err := createOneTcpConn(address, connTimeout, maxRetry)
		if err != nil {
			cp.M <- nil
			continue
		}
		cp.M <- conn
	}
	return cp
}

func createOneTcpConn(address string, connTimeout, retry int) (conn net.Conn, err error) {
	ct := time.Duration(connTimeout) * time.Millisecond

	for i := 0; i < retry; i++ {
		_, err = net.ResolveTCPAddr("tcp", address)
		if err != nil {
			continue
		}

		conn, err = net.DialTimeout("tcp", address, ct)
		if err == nil {
			break
		}
	}

	if err != nil {
		log.Printf("build %s tcp connect fail: %v", address, err)
	}

	return
}
