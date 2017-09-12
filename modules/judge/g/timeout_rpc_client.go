package g

import (
	"time"
	"errors"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type TimeoutRpcClient struct {
	*rpc.Client
}

func NewTimeoutRpcClient(network, address string, timeout time.Duration) (*TimeoutRpcClient, error) {
	conn, err := net.DialTimeout(network, address, timeout)
	if err != nil {
		return nil, err
	}

	return &TimeoutRpcClient {
		Client: jsonrpc.NewClient(conn),
	}, nil
}

func (this *TimeoutRpcClient) CallTimeout(serviceMethod string, args interface{}, reply interface{}, timeout time.Duration) error {
	call := this.Go(serviceMethod, args, reply, make(chan *rpc.Call, 1))

	select {
	case callPtr := <-call.Done:
		return callPtr.Error
	case <-time.After(timeout):
		return errors.New("rpc call timeout")
	}
}
