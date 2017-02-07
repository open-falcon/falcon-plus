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

func (this TsdbClient) Name() string {
	return this.name
}

func (this TsdbClient) Closed() bool {
	return this.cli == nil
}

func (this TsdbClient) Close() error {
	if this.cli != nil {
		err := this.cli.Close()
		this.cli = nil
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

func (this *TsdbConnPoolHelper) Send(data []byte) (err error) {
	conn, err := this.p.Fetch()
	if err != nil {
		return fmt.Errorf("get connection fail: err %v. proc: %s", err, this.p.Proc())
	}

	cli := conn.(TsdbClient).cli

	done := make(chan error, 1)
	go func() {
		_, err = cli.Write(data)
		done <- err
	}()

	select {
	case <-time.After(time.Duration(this.callTimeout) * time.Millisecond):
		this.p.ForceClose(conn)
		return fmt.Errorf("%s, call timeout", this.address)
	case err = <-done:
		if err != nil {
			this.p.ForceClose(conn)
			err = fmt.Errorf("%s, call failed, err %v. proc: %s", this.address, err, this.p.Proc())
		} else {
			this.p.Release(conn)
		}
		return err
	}

	return
}

func (this *TsdbConnPoolHelper) Destroy() {
	if this.p != nil {
		this.p.Destroy()
	}
}
