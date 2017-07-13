package goperfcounter

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/niean/go-metrics-lite"
	bhttp "github.com/niean/gotools/http/httpclient/beego"
)

const (
	GAUGE = "GAUGE"
)

func pushToFalcon() {
	cfg := config()
	step := cfg.Step
	api := cfg.Push.Api
	debug := cfg.Debug

	// align push start ts
	alignPushStartTs(step)

	for _ = range time.Tick(time.Duration(step) * time.Second) {
		selfMeter("pfc.push.cnt", 1) // statistics

		fms := falconMetrics()
		start := time.Now()
		err := push(fms, api, debug)
		selfGauge("pfc.push.ms", int64(time.Since(start)/time.Millisecond)) // statistics

		if err != nil {
			if debug {
				log.Printf("[perfcounter] send to %s error: %v", api, err)
			}
			selfGauge("pfc.push.size", int64(0)) // statistics
		} else {
			selfGauge("pfc.push.size", int64(len(fms))) // statistics
		}
	}
}

func falconMetric(types []string) []*MetricValue {
	fd := []*MetricValue{}
	for _, ty := range types {
		if r, ok := values[ty]; ok && r != nil {
			data := _falconMetric(r)
			fd = append(fd, data...)
		}
	}
	return fd
}

func falconMetrics() []*MetricValue {
	data := make([]*MetricValue, 0)
	for _, r := range values {
		nd := _falconMetric(r)
		data = append(data, nd...)
	}
	return data
}

// internal
func _falconMetric(r metrics.Registry) []*MetricValue {
	cfg := config()
	endpoint := cfg.Hostname
	step := cfg.Step
	tags := cfg.Tags
	ts := time.Now().Unix()

	data := make([]*MetricValue, 0)
	r.Each(func(name string, i interface{}) {
		switch metric := i.(type) {
		case metrics.Gauge:
			m := gaugeMetricValue(metric, name, endpoint, tags, step, ts)
			data = append(data, m...)
		case metrics.GaugeFloat64:
			m := gaugeFloat64MetricValue(metric, name, endpoint, tags, step, ts)
			data = append(data, m...)
		case metrics.Counter:
			m := counterMetricValue(metric, name, endpoint, tags, step, ts)
			data = append(data, m...)
		case metrics.Meter:
			m := metric.Snapshot()
			ms := meterMetricValue(m, name, endpoint, tags, step, ts)
			data = append(data, ms...)
		case metrics.Histogram:
			h := metric.Snapshot()
			ms := histogramMetricValue(h, name, endpoint, tags, step, ts)
			data = append(data, ms...)
		}
	})

	return data
}

func gaugeMetricValue(metric metrics.Gauge, metricName, endpoint, oldtags string, step, ts int64) []*MetricValue {
	tags := getTags(metricName, oldtags)
	c := newMetricValue(endpoint, "value", metric.Value(), step, GAUGE, tags, ts)
	return []*MetricValue{c}
}

func gaugeFloat64MetricValue(metric metrics.GaugeFloat64, metricName, endpoint, oldtags string, step, ts int64) []*MetricValue {
	tags := getTags(metricName, oldtags)
	c := newMetricValue(endpoint, "value", metric.Value(), step, GAUGE, tags, ts)
	return []*MetricValue{c}
}

func counterMetricValue(metric metrics.Counter, metricName, endpoint, oldtags string, step, ts int64) []*MetricValue {
	tags := getTags(metricName, oldtags)
	c1 := newMetricValue(endpoint, "count", metric.Count(), step, GAUGE, tags, ts)
	return []*MetricValue{c1}
}

func meterMetricValue(metric metrics.Meter, metricName, endpoint, oldtags string, step, ts int64) []*MetricValue {
	data := make([]*MetricValue, 0)
	tags := getTags(metricName, oldtags)

	c1 := newMetricValue(endpoint, "rate", metric.RateStep(), step, GAUGE, tags, ts)
	c2 := newMetricValue(endpoint, "sum", metric.Count(), step, GAUGE, tags, ts)
	data = append(data, c1, c2)

	return data
}

func histogramMetricValue(metric metrics.Histogram, metricName, endpoint, oldtags string, step, ts int64) []*MetricValue {
	data := make([]*MetricValue, 0)
	tags := getTags(metricName, oldtags)

	values := make(map[string]interface{})
	ps := metric.Percentiles([]float64{0.75, 0.95, 0.99})
	values["min"] = metric.Min()
	values["max"] = metric.Max()
	values["mean"] = metric.Mean()
	values["75th"] = ps[0]
	values["95th"] = ps[1]
	values["99th"] = ps[2]
	for key, val := range values {
		c := newMetricValue(endpoint, key, val, step, GAUGE, tags, ts)
		data = append(data, c)
	}

	return data
}

func newMetricValue(endpoint, metric string, value interface{}, step int64, t, tags string, ts int64) *MetricValue {
	return &MetricValue{
		Endpoint:  endpoint,
		Metric:    metric,
		Value:     value,
		Step:      step,
		Type:      t,
		Tags:      tags,
		Timestamp: ts,
	}
}

func getTags(name string, tags string) string {
	if tags == "" {
		return fmt.Sprintf("name=%s", name)
	}
	return fmt.Sprintf("%s,name=%s", tags, name)
}

//
func push(data []*MetricValue, url string, debug bool) error {
	dlen := len(data)
	pkg := 200 //send pkg items once
	sent := 0
	for {
		if sent >= dlen {
			break
		}

		end := sent + pkg
		if end > dlen {
			end = dlen
		}

		pkgData := data[sent:end]
		jr, err := json.Marshal(pkgData)
		if err != nil {
			return err
		}

		response, err := bhttp.Post(url).Body(jr).String()
		if err != nil {
			return err
		}
		sent = end

		if debug {
			log.Printf("[perfcounter] push result: %v, data: %v\n", response, pkgData)
		}
	}
	return nil
}

//
func alignPushStartTs(stepSec int64) {
	nw := time.Duration(time.Now().UnixNano())
	step := time.Duration(stepSec) * time.Second
	sleepNano := step - nw%step
	if sleepNano > 0 {
		time.Sleep(sleepNano)
	}
}

//
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
	return fmt.Sprintf(
		"<Endpoint:%s, Metric:%s, Tags:%s, Type:%s, Step:%d, Timestamp:%d, Value:%v>",
		this.Endpoint,
		this.Metric,
		this.Tags,
		this.Type,
		this.Step,
		this.Timestamp,
		this.Value,
	)
}
