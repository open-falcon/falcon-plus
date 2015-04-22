package cron

import (
	"github.com/open-falcon/agent/g"
	"github.com/open-falcon/common/model"
	"log"
	"strconv"
	"strings"
	"time"
)

func SyncBuiltinItems() {
	if g.Config().Heartbeat.Enabled && g.Config().Heartbeat.Addr != "" {
		go syncBuiltinItems()
	}
}

func syncBuiltinItems() {

	var timestamp int64
	var checksum string

	duration := time.Duration(g.Config().Heartbeat.Interval) * time.Second

	for {
	REST:
		time.Sleep(duration)

		var ports = []int64{}
		var procs = make(map[string]map[int]string)

		hostname, err := g.Hostname()
		if err != nil {
			goto REST
		}

		req := model.AgentHeartbeatRequest{
			Hostname: hostname,
			Checksum: checksum,
		}

		var resp g.BuiltinItemResp
		err = g.HbsClient.Call("Agent.GetBuiltinItems", req, &resp)
		if err != nil {
			log.Println("[ERROR]", err)
			goto REST
		}

		if resp.Timestamp <= timestamp {
			log.Println("resp.Timestamp <= timestamp")
			goto REST
		}

		if resp.Checksum == checksum {
			goto REST
		}

		timestamp = resp.Timestamp
		checksum = resp.Checksum

		for _, item := range resp.Items {
			if item.Metric == "net.port.listen" {
				if port, err := strconv.ParseInt(item.Tags[5:], 10, 64); err == nil {
					ports = append(ports, port)
				}

				continue
			}

			if item.Metric == "proc.num" {
				arr := strings.Split(item.Tags, ",")

				tmpMap := make(map[int]string)

				for i := 0; i < len(arr); i++ {
					if strings.HasPrefix(arr[i], "name=") {
						tmpMap[1] = arr[i][5:]
					} else if strings.HasPrefix(arr[i], "cmdline=") {
						tmpMap[2] = arr[i][8:]
					}
				}

				procs[item.Tags] = tmpMap
			}
		}

		g.SetReportPorts(ports)
		g.SetReportProcs(procs)

	}
}
