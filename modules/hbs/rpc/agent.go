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

package rpc

import (
	"bytes"
	"github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/hbs/cache"
	"github.com/open-falcon/falcon-plus/modules/hbs/g"
	"sort"
	"strings"
	"time"
)

func (t *Agent) MinePlugins(args model.AgentHeartbeatRequest, reply *model.AgentPluginsResponse) error {
	if args.Hostname == "" {
		return nil
	}

	reply.Plugins = cache.GetPlugins(args.Hostname)
	reply.Timestamp = time.Now().Unix()

	return nil
}

func (t *Agent) ReportStatus(args *model.AgentReportRequest, reply *model.SimpleRpcResponse) error {
	if args.Hostname == "" {
		reply.Code = 1
		return nil
	}

	cache.Agents.Put(args)

	return nil
}

// 需要checksum一下来减少网络开销？其实白名单通常只会有一个或者没有，无需checksum
func (t *Agent) TrustableIps(args *model.NullRpcRequest, ips *string) error {
	*ips = strings.Join(g.Config().Trustable, ",")
	return nil
}

// agent按照server端的配置，按需采集的metric，比如net.port.listen port=22 或者 proc.num name=zabbix_agentd
func (t *Agent) BuiltinMetrics(args *model.AgentHeartbeatRequest, reply *model.BuiltinMetricResponse) error {
	if args.Hostname == "" {
		return nil
	}

	metrics, err := cache.GetBuiltinMetrics(args.Hostname)
	if err != nil {
		return nil
	}

	checksum := ""
	if len(metrics) > 0 {
		checksum = DigestBuiltinMetrics(metrics)
	}

	if args.Checksum == checksum {
		reply.Metrics = []*model.BuiltinMetric{}
	} else {
		reply.Metrics = metrics
	}
	reply.Checksum = checksum
	reply.Timestamp = time.Now().Unix()

	return nil
}

func DigestBuiltinMetrics(items []*model.BuiltinMetric) string {
	sort.Sort(model.BuiltinMetricSlice(items))

	var buf bytes.Buffer
	for _, m := range items {
		buf.WriteString(m.String())
	}

	return utils.Md5(buf.String())
}
