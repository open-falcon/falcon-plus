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

package model

import (
	"fmt"
)

type AgentReportRequest struct {
	Hostname      string
	IP            string
	AgentVersion  string
	PluginVersion string
}

func (this *AgentReportRequest) String() string {
	return fmt.Sprintf(
		"<Hostname:%s, IP:%s, AgentVersion:%s, PluginVersion:%s>",
		this.Hostname,
		this.IP,
		this.AgentVersion,
		this.PluginVersion,
	)
}

type AgentUpdateInfo struct {
	LastUpdate    int64
	ReportRequest *AgentReportRequest
}

type AgentHeartbeatRequest struct {
	Hostname string
	Checksum string
}

func (this *AgentHeartbeatRequest) String() string {
	return fmt.Sprintf(
		"<Hostname: %s, Checksum: %s>",
		this.Hostname,
		this.Checksum,
	)
}

type AgentPluginsResponse struct {
	Plugins   []string
	Timestamp int64
}

func (this *AgentPluginsResponse) String() string {
	return fmt.Sprintf(
		"<Plugins:%v, Timestamp:%v>",
		this.Plugins,
		this.Timestamp,
	)
}

// e.g. net.port.listen or proc.num
type BuiltinMetric struct {
	Metric string
	Tags   string
}

func (this *BuiltinMetric) String() string {
	return fmt.Sprintf(
		"%s/%s",
		this.Metric,
		this.Tags,
	)
}

type BuiltinMetricResponse struct {
	Metrics   []*BuiltinMetric
	Checksum  string
	Timestamp int64
}

func (this *BuiltinMetricResponse) String() string {
	return fmt.Sprintf(
		"<Metrics:%v, Checksum:%s, Timestamp:%v>",
		this.Metrics,
		this.Checksum,
		this.Timestamp,
	)
}

type BuiltinMetricSlice []*BuiltinMetric

func (this BuiltinMetricSlice) Len() int {
	return len(this)
}
func (this BuiltinMetricSlice) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}
func (this BuiltinMetricSlice) Less(i, j int) bool {
	return this[i].String() < this[j].String()
}
