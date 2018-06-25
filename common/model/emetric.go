package model

type EMetric struct {
	Endpoint  string                 `json:"endpoint"`
	Key       string                 `json:"key"`
	Filters   map[string]interface{} `json:"values"`
	Timestamp int64                  `json:"timestamp"`
}

func NewEMetric() *EMetric {
	e := EMetric{}
	e.Filters = map[string]interface{}{}
	return &e
}

func (e *EMetric) PK() string {
	return e.Key
}
