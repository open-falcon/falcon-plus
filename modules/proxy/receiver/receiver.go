package receiver

import (
	"github.com/open-falcon/falcon-plus/modules/proxy/receiver/rpc"
	"github.com/open-falcon/falcon-plus/modules/proxy/receiver/socket"
)

func Start() {
	go rpc.StartRpc()
	go socket.StartSocket()
}
