package g

import (
	"github.com/toolkits/net"
	"log"
)

var LocalIps []string

func InitVars() {
	var err error
	LocalIps, err = net.IntranetIP()
	if err != nil {
		log.Fatalln("get intranet ip fail:", err)
	}
}
