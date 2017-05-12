package backend_pool

import (
	"fmt"
	"net"
	"time"

	connp "github.com/toolkits/conn_pool"
)

// TSDB
type TsdbClient struct {
	cli  net.Conn
	name string
}

func (t TsdbClient) Name() string {
	return t.name
}

func (t TsdbClient) Closed() bool {
	return t.cli == nil
}

func (t TsdbClient) Close() error {
	if t.cli != nil {
		err := t.cli.Close()
		t.cli = nil
		return err
	}
	return nil
}

func newTsdbConnPool(address string, maxConns int, maxIdle int, connTimeout int) *connp.ConnPool {
	pool := connp.NewConnPool("tsdb", address, int32(maxConns), int32(maxIdle))

	pool.New = func(name string) (connp.NConn, error) {
		_, err := net.ResolveTCPAddr("tcp", address)
		if err != nil {
			return nil, err
		}

		conn, err := net.DialTimeout("tcp", address, time.Duration(connTimeout)*time.Millisecond)
		if err != nil {
			return nil, err
		}

		return TsdbClient{conn, name}, nil
	}

	return pool
}

type TsdbConnPoolHelper struct {
	p           *connp.ConnPool
	maxConns    int
	maxIdle     int
	connTimeout int
	callTimeout int
	address     string
}

func NewTsdbConnPoolHelper(address string, maxConns, maxIdle, connTimeout, callTimeout int) *TsdbConnPoolHelper {
	return &TsdbConnPoolHelper{
		p:           newTsdbConnPool(address, maxConns, maxIdle, connTimeout),
		maxConns:    maxConns,
		maxIdle:     maxIdle,
		connTimeout: connTimeout,
		callTimeout: callTimeout,
		address:     address,
	}
}

func (t *TsdbConnPoolHelper) Send(data []byte) (err error) {
	conn, err := t.p.Fetch()
	if err != nil {
		return fmt.Errorf("get connection fail: err %v. proc: %s", err, t.p.Proc())
	}

	cli := conn.(TsdbClient).cli

	done := make(chan error, 1)
	go func() {
		_, err = cli.Write(data)
		done <- err
	}()

	select {
	case <-time.After(time.Duration(t.callTimeout) * time.Millisecond):
		t.p.ForceClose(conn)
		return fmt.Errorf("%s, call timeout", t.address)
	case err = <-done:
		if err != nil {
			t.p.ForceClose(conn)
			err = fmt.Errorf("%s, call failed, err %v. proc: %s", t.address, err, t.p.Proc())
		} else {
			t.p.Release(conn)
		}
		return err
	}
}

func (t *TsdbConnPoolHelper) Destroy() {
	if t.p != nil {
		t.p.Destroy()
	}
}
