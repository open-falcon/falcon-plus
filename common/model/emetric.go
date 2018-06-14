package model

type EMetric struct {
	Endpoint  string             `json:"endpoint"`
	Metric    string             `json:"metric"`
	Values    map[string]float64 `json:"values"`
	Filters   map[string]string  `json:"filters"`
	Timestamp int64              `json:"timestamp"`
}

func NewEMetric() *EMetric {
	e := EMetric{}
	e.Values = map[string]float64{}
	e.Filters = map[string]string{}
	return &e
}

func (e *EMetric) PK() string {
	return e.Metric
}
