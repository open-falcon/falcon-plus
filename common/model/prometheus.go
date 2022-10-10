package model

// MetricType Prometheus指标类型：GAUGE|COUNTER
type P8sItem struct {
	Endpoint   string            `json:"endpoint"`
	Metric     string            `json:"metric"`
	Tags       map[string]string `json:"tags"`
	Value      float64           `json:"value"`
	Timestamp  int64             `json:"timestamp"`
	MetricType string            `json:"metric_type"`
	Step       int               `json:"step"`
}
