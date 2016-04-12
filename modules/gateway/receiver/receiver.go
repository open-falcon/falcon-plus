package receiver

import (
	"github.com/open-falcon/gateway/receiver/rpc"
	"github.com/open-falcon/gateway/receiver/socket"
)

func Start() {
	go rpc.StartRpc()
	go socket.StartSocket()
}
