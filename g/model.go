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
	Msg     string
	Total   int
	Latency int64
}

func (this *TransferResp) String() string {
	return fmt.Sprintf("<Total=%v, Latency=%vms, Msg:%s>",
		this.Total,
		this.Latency,
		this.Msg)
}
