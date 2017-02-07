package receiver

import (
	"github.com/open-falcon/falcon-plus/modules/transfer/receiver/rpc"
	"github.com/open-falcon/falcon-plus/modules/transfer/receiver/socket"
)

func Start() {
	go rpc.StartRpc()
	go socket.StartSocket()
}
