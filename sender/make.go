package sender

import (
	"time"

	"github.com/open-falcon/common/model"
)

func MakeMetaData(endpoint, metric, tags string, val interface{}, counterType string, step_and_ts ...int64) *model.JsonMetaData {
	md := model.JsonMetaData{
		Endpoint:    endpoint,
		Metric:      metric,
		Tags:        tags,
		Value:       val,
		CounterType: counterType,
	}

	argc := len(step_and_ts)
	if argc == 0 {
		md.Step = 60
		md.Timestamp = time.Now().Unix()
	} else if argc == 1 {
		md.Step = step_and_ts[0]
		md.Timestamp = time.Now().Unix()
	} else if argc == 2 {
		md.Step = step_and_ts[0]
		md.Timestamp = step_and_ts[1]
	}

	return &md
}

func MakeGaugeValue(endpoint, metric, tags string, val interface{}, step_and_ts ...int64) *model.JsonMetaData {
	return MakeMetaData(endpoint, metric, tags, val, "GAUGE", step_and_ts...)
}

func MakeCounterValue(endpoint, metric, tags string, val interface{}, step_and_ts ...int64) *model.JsonMetaData {
	return MakeMetaData(endpoint, metric, tags, val, "COUNTER", step_and_ts...)
}

func PushGauge(endpoint, metric, tags string, val interface{}, step_and_ts ...int64) {
	md := MakeGaugeValue(endpoint, metric, tags, val, step_and_ts...)
	MetaDataQueue.PushFront(md)
}

func PushCounter(endpoint, metric, tags string, val interface{}, step_and_ts ...int64) {
	md := MakeCounterValue(endpoint, metric, tags, val, step_and_ts...)
	MetaDataQueue.PushFront(md)
}

func Push(endpoint, metric, tags string, val interface{}, counterType string, step_and_ts ...int64) {
	md := MakeMetaData(endpoint, metric, tags, val, counterType, step_and_ts...)
	MetaDataQueue.PushFront(md)
}
