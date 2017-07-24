package goperfcounter

import (
	"time"

	"github.com/niean/go-metrics-lite"
	"github.com/niean/goperfcounter/base"
)

var (
	gpGaugeFloat64 = metrics.NewRegistry()
	gpCounter      = metrics.NewRegistry()
	gpMeter        = metrics.NewRegistry()
	gpHistogram    = metrics.NewRegistry()
	gpDebug        = metrics.NewRegistry()
	gpRuntime      = metrics.NewRegistry()
	gpSelf         = metrics.NewRegistry()
	values         = make(map[string]metrics.Registry) //readonly,mappings of metrics
)

func init() {
	values["gauge"] = gpGaugeFloat64
	values["counter"] = gpCounter
	values["meter"] = gpMeter
	values["histogram"] = gpHistogram
	values["debug"] = gpDebug
	values["runtime"] = gpRuntime
	values["self"] = gpSelf
}

//
func rawMetric(types []string) map[string]interface{} {
	data := make(map[string]interface{})
	for _, mtype := range types {
		if v, ok := values[mtype]; ok {
			data[mtype] = v.Values()
		}
	}
	return data
}

func rawMetrics() map[string]interface{} {
	data := make(map[string]interface{})
	for key, v := range values {
		data[key] = v.Values()
	}
	return data
}

func rawSizes() map[string]int64 {
	data := map[string]int64{}
	all := int64(0)
	for key, v := range values {
		kv := v.Size()
		all += kv
		data[key] = kv
	}
	data["all"] = all
	return data
}

func collectBase(bases []string) {
	// start base collect after 30sec
	time.Sleep(time.Duration(30) * time.Second)

	if contains(bases, "debug") {
		base.RegisterAndCaptureDebugGCStats(gpDebug, 5e9)
	}

	if contains(bases, "runtime") {
		base.RegisterAndCaptureRuntimeMemStats(gpRuntime, 5e9)
	}
}

func contains(bases []string, name string) bool {
	for _, n := range bases {
		if n == name {
			return true
		}
	}
	return false
}
