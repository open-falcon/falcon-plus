// Copyright 2017 Xiaomi, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package g

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/open-falcon/falcon-plus/common/model"
	"github.com/toolkits/slice"
)

var Root string

func InitRootDir() {
	var err error
	Root, err = os.Getwd()
	if err != nil {
		log.Fatalln("getwd fail:", err)
	}
}

var LocalIp string

func InitLocalIp() {
	if Config().Heartbeat.Enabled {
		for _, addr := range Config().Heartbeat.Addrs {
			conn, err := net.DialTimeout("tcp", addr, time.Second*10)
			if err != nil {
				log.Println(fmt.Sprintf("connect to heartbeat server %s failed", addr))
			} else {
				defer conn.Close()
				LocalIp = strings.Split(conn.LocalAddr().String(), ":")[0]
				break
			}
		}
	} else {
		log.Println("hearbeat is not enabled, can't get localip")
	}
}

var (
	HbsClient *SingleConnRpcClient
)

func InitRpcClients() {
	if Config().Heartbeat.Enabled && len(Config().Heartbeat.Addrs) != 0 {
		HbsClient = &SingleConnRpcClient{
			RpcServers: Config().Heartbeat.Addrs,
			Timeout:    time.Duration(Config().Heartbeat.Timeout) * time.Millisecond,
		}
	}

}

func SendToTransfer(metrics []*model.MetricValue) {
	if len(metrics) == 0 {
		return
	}

	dt := Config().DefaultTags
	if len(dt) > 0 {
		var buf bytes.Buffer
		default_tags_list := []string{}
		for k, v := range dt {
			buf.Reset()
			buf.WriteString(k)
			buf.WriteString("=")
			buf.WriteString(v)
			default_tags_list = append(default_tags_list, buf.String())
		}
		default_tags := strings.Join(default_tags_list, ",")

		for i, x := range metrics {
			buf.Reset()
			if x.Tags == "" {
				metrics[i].Tags = default_tags
			} else {
				buf.WriteString(metrics[i].Tags)
				buf.WriteString(",")
				buf.WriteString(default_tags)
				metrics[i].Tags = buf.String()
			}
		}
	}

	debug := Config().Debug

	if debug {
		log.Printf("=> <Total=%d> %v\n", len(metrics), metrics[0])
	}

	var resp model.TransferResponse
	SendMetrics(metrics, &resp)

	if debug {
		log.Println("<=", &resp)
	}
}

var (
	reportUrls     map[string]string
	reportUrlsLock = new(sync.RWMutex)
)

func ReportUrls() map[string]string {
	reportUrlsLock.RLock()
	defer reportUrlsLock.RUnlock()
	return reportUrls
}

func SetReportUrls(urls map[string]string) {
	reportUrlsLock.RLock()
	defer reportUrlsLock.RUnlock()
	reportUrls = urls
}

var (
	reportPorts     []int64
	reportPortsLock = new(sync.RWMutex)
)

func ReportPorts() []int64 {
	reportPortsLock.RLock()
	defer reportPortsLock.RUnlock()
	return reportPorts
}

func SetReportPorts(ports []int64) {
	reportPortsLock.Lock()
	defer reportPortsLock.Unlock()
	reportPorts = ports
}

var (
	duPaths     []string
	duPathsLock = new(sync.RWMutex)
)

func DuPaths() []string {
	duPathsLock.RLock()
	defer duPathsLock.RUnlock()
	return duPaths
}

func SetDuPaths(paths []string) {
	duPathsLock.Lock()
	defer duPathsLock.Unlock()
	duPaths = paths
}

var (
	// tags => {1=>name, 2=>cmdline}
	// e.g. 'name=falcon-agent'=>{1=>falcon-agent}
	// e.g. 'cmdline=xx'=>{2=>xx}
	reportProcs     map[string]map[int]string
	reportProcsLock = new(sync.RWMutex)
)

func ReportProcs() map[string]map[int]string {
	reportProcsLock.RLock()
	defer reportProcsLock.RUnlock()
	return reportProcs
}

func SetReportProcs(procs map[string]map[int]string) {
	reportProcsLock.Lock()
	defer reportProcsLock.Unlock()
	reportProcs = procs
}

var (
	ips     []string
	ipsLock = new(sync.Mutex)
)

func TrustableIps() []string {
	ipsLock.Lock()
	defer ipsLock.Unlock()
	return ips
}

func SetTrustableIps(ipStr string) {
	arr := strings.Split(ipStr, ",")
	ipsLock.Lock()
	defer ipsLock.Unlock()
	ips = arr
}

func IsTrustable(remoteAddr string) bool {
	ip := remoteAddr
	idx := strings.LastIndex(remoteAddr, ":")
	if idx > 0 {
		ip = remoteAddr[0:idx]
	}

	if ip == "127.0.0.1" {
		return true
	}

	return slice.ContainsString(TrustableIps(), ip)
}
