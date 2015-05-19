package receiver

import (
	"github.com/open-falcon/transfer/receiver/rpc"
	"github.com/open-falcon/transfer/receiver/socket"
)

func Start() {
	go rpc.StartRpc()
	go socket.StartSocket()
}
