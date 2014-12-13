package cron

import (
	"github.com/open-falcon/agent/g"
	"log"
	"os"
	"sync"
	"time"
)

var (
	whiteIPs    []*g.WhiteIP
	whiteIpLock = new(sync.Mutex)
)

func WhiteIps() []*g.WhiteIP {
	whiteIpLock.Lock()
	defer whiteIpLock.Unlock()
	return whiteIPs
}

func setWhiteIps(ips []*g.WhiteIP) {
	whiteIpLock.Lock()
	defer whiteIpLock.Unlock()
	whiteIPs = ips
}

func SyncWhiteIPs() {
	if g.Config().Heartbeat.Enabled && g.Config().Heartbeat.Addr != "" {
		go syncWhiteIPs()
	}
}

func syncWhiteIPs() {
	var ipsChecksum string
	var ipsTimestamp int64

	duration := time.Duration(g.Config().Heartbeat.Interval) * time.Second

	for {
	REST:
		time.Sleep(duration)

		hostname, err := os.Hostname()
		if err != nil {
			log.Println("[ERROR] os.Hostname() fail:", err)
			goto REST
		}

		req := g.AgentReq{
			Host:     g.Host{HostName: hostname},
			Checksum: ipsChecksum,
		}

		var resp g.IpWhiteListResp

		err = g.HbsClient.Call("Agent.GetWhiteIPList", req, &resp)
		if err != nil {
			log.Println("[ERROR]", err)
			goto REST
		}

		if resp.Checksum == "" {
			log.Println("[ERROR] resp.Checksum is blank")
			goto REST
		}

		if resp.Timestamp <= ipsTimestamp {
			log.Println("resp.Timestamp <= ipsTimestamp")
			goto REST
		}

		if resp.Checksum == ipsChecksum {
			goto REST
		}

		setWhiteIps(resp.Ips)
		ipsTimestamp = resp.Timestamp
		ipsChecksum = resp.Checksum
	}
}
