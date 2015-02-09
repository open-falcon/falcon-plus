package g

import (
	"fmt"
)

type MetricValue struct {
	Endpoint  string      `json:"endpoint"`
	Metric    string      `json:"metric"`
	Value     interface{} `json:"value"`
	Step      int64       `json:"step"`
	Type      string      `json:"counterType"`
	Tags      string      `json:"tags"`
	Timestamp int64       `json:"timestamp"`
}

func (this *MetricValue) String() string {
	return fmt.Sprintf("<Endpoint:%s, Metric:%s, Type:%s, Tags:%s, Step:%d, Time:%d, Value:%v>",
		this.Endpoint,
		this.Metric,
		this.Type,
		this.Tags,
		this.Step,
		this.Timestamp,
		this.Value)
}

type TransferResp struct {
	Msg        string
	Total      int
	ErrInvalid int
	Latency    int64
}

func (this *TransferResp) String() string {
	return fmt.Sprintf("<Total=%v, Latency=%vms, Invalid:%v, Msg:%s>",
		this.Total,
		this.Latency,
		this.ErrInvalid,
		this.Msg)
}

type AgentReportReq struct {
	HostName      string
	Version       string
	Meta          string
	PluginVersion string
}

type AgentReportResp struct {
	Status bool
	Msg    string
}

type Plugin struct {
	Path string
}

type Host struct {
	HostId   int
	HostName string
	Pause    int
	Uuid     string
}

type AgentReq struct {
	Host
	Checksum string
}

type AgentPluginsResp struct {
	Plugins   []*Plugin
	HostName  string
	Checksum  string
	Timestamp int64
}

type WhiteIP struct {
	Ip string
}

type IpWhiteListResp struct {
	Ips       []*WhiteIP
	Checksum  string
	Timestamp int64
}

type BuiltinItem struct {
	Metric string
	Tags   string
}

type BuiltinItemResp struct {
	Items     []*BuiltinItem
	Checksum  string
	Timestamp int64
}
