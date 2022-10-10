package rpc

// code == 0 => success
// code == 1 => bad request
type SimpleRpcResponse struct {
	Code int `json:"code"`
}

type NullRpcRequest struct {
}

type P8sItem struct {
	Endpoint       string            `json:"endpoint"`
	Metric         string            `json:"metric"`
	Tags           map[string]string `json:"tags"`
	TagKeys        []string          `json:"-"`
	TagValues      []string          `json:"-"`
	Value          float64           `json:"value"`
	Timestamp      int64             `json:"timestamp"`
	MetricType     string            `json:"metric_type"`
	Step           int               `json:"step"`
	PK             string            `json:"-"`
	PKWithTagValue string            `json:"-"`
}
