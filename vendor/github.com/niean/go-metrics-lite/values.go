package metrics

func (r *StandardRegistry) Values() interface{} {
	data := make(map[string]map[string]interface{})
	r.Each(func(name string, i interface{}) {
		values := make(map[string]interface{})
		switch metric := i.(type) {
		case Counter:
			values["count"] = metric.Count()
		case Gauge:
			values["value"] = metric.Value()
		case GaugeFloat64:
			values["value"] = metric.Value()
		case Healthcheck:
			values["error"] = nil
			metric.Check()
			if err := metric.Error(); nil != err {
				values["error"] = metric.Error().Error()
			}
		case Meter:
			m := metric.Snapshot()
			values["sum"] = m.Count()
			values["rate"] = m.RateStep()
			values["rate.1min"] = m.Rate1()
			values["rate.5min"] = m.Rate5()
			values["rate.15min"] = m.Rate15()
		case Histogram:
			h := metric.Snapshot()
			ps := h.Percentiles([]float64{0.75, 0.95, 0.99})
			values["min"] = h.Min()
			values["max"] = h.Max()
			values["mean"] = h.Mean()
			values["75th"] = ps[0]
			values["95th"] = ps[1]
			values["99th"] = ps[2]
		case Timer:
			t := metric.Snapshot()
			ps := t.Percentiles([]float64{0.75, 0.95, 0.99})
			values["min"] = t.Min()
			values["max"] = t.Max()
			values["mean"] = t.Mean()
			values["75th"] = ps[0]
			values["95th"] = ps[1]
			values["99th"] = ps[2]
			values["sum"] = t.Count()
			values["rate"] = t.RateStep()
			values["rate.1min"] = t.Rate1()
			values["rate.5min"] = t.Rate5()
			values["rate.15min"] = t.Rate15()
		}
		data[name] = values
	})
	return data
}
