package graph

import (
	"errors"
	"github.com/open-falcon/common/model"
	"github.com/open-falcon/query/g"
	"github.com/toolkits/cpool"
	"github.com/toolkits/logger"
	"github.com/toolkits/rpool/conn_pool"
	"io"
	"net"
	"net/rpc"
	"time"
)

var (
	backend *cpool.RingBackend = &cpool.RingBackend{
		Addrs: map[string][]string{},
		Ring:  nil,
		Pools: make(map[string]*conn_pool.ConnPool),
	}
)

type RpcConn struct {
	cli *rpc.Client
}

func (t RpcConn) Close() error {
	if t.cli == nil {
		return errors.New("nil conn")
	}
	err := t.cli.Close()
	t.cli = nil
	return err
}

func ReloadBackends() {
	cfg := g.Config().Graph
	for {
		time.Sleep(time.Duration(cfg.ReloadInterval) * time.Second)
		err := InitBackends()
		if err != nil {
			logger.Error("reload backends fail: %v", err)
		}
	}
}

func InitBackends() error {
	var err error
	cfg := g.Config().Graph

	file_path := cfg.Backends
	err = backend.LoadAddrs(file_path)
	if err != nil {
		return err
	}

	backend.InitRing(cfg.Replicas)

	err = initConnPools()
	if err != nil {
		return err
	}

	return nil
}

func DestroyConnPools() {
	backend.DestroyConnPools()
}

func initConnPools() error {
	cfg := g.Config()
	if cfg.LogLevel == "trace" || cfg.LogLevel == "debug" {
		conn_pool.EnableSlowLog(true, cfg.SlowLog)
	}

	var (
		tmp_addrs map[string][]string
		tmp_pools map[string]*conn_pool.ConnPool
	)

	backend.RLock()
	tmp_addrs = backend.Addrs
	tmp_pools = backend.Pools
	backend.RUnlock()

	c := cfg.Graph
	for name, addr_list := range tmp_addrs {
		for _, addr := range addr_list {
			if _, ok := tmp_pools[addr]; !ok {
				pool := conn_pool.NewConnPool(addr, c.MaxConns, c.MaxIdle)

				pool.New = func() (io.Closer, error) {
					_, err := net.ResolveTCPAddr("tcp", pool.Name)
					if err != nil {
						return nil, err
					}
					conn, err := net.DialTimeout("tcp", pool.Name, time.Duration(c.Timeout)*time.Millisecond)
					if err != nil {
						return nil, err
					}

					return RpcConn{rpc.NewClient(conn)}, nil
				}

				pool.Ping = func(conn io.Closer) error {
					rpc_conn := conn.(RpcConn)
					if rpc_conn.cli == nil {
						return errors.New("nil conn")
					}

					resp := &model.SimpleRpcResponse{}
					err := rpc_conn.cli.Call("Graph.Ping", model.NullRpcRequest{}, resp)
					logger.Trace("Graph.Ping resp: %v", resp)

					return err
				}

				tmp_pools[addr] = pool
				logger.Info("create the pool: %s %s", name, addr)
			} else {
				logger.Trace("keep the pool: %s %s", name, addr)
			}
		}
	}

	backend.Lock()
	defer backend.Unlock()
	backend.Pools = tmp_pools

	return nil
}
