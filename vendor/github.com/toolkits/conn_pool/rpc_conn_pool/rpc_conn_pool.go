package rpc_conn_pool

import (
	"net/rpc"
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

func NewRpcClient(cli *rpc.Client, name string) *RpcClient {
	return &RpcClient{cli: cli, name: name}
}

func NewRpcClientWithCodec(codec rpc.ClientCodec, name string) *RpcClient {
	return &RpcClient{
		cli:  rpc.NewClientWithCodec(codec),
		name: name,
	}
}
