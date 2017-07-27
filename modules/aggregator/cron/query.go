package cron

import (
	"github.com/open-falcon/falcon-plus/modules/aggregator/sdk"
)

func queryCounterLast(numeratorOperands, denominatorOperands, hostnames []string, begin, end int64) (map[string]float64, error) {
	counters := []string{}

	counters = append(counters, numeratorOperands...)
	counters = append(counters, denominatorOperands...)

	resp, err := sdk.QueryLastPoints(hostnames, counters)
	if err != nil {
		return map[string]float64{}, err
	}

	ret := make(map[string]float64)
	for _, res := range resp {
		v := res.Value
		if v.Timestamp < begin || v.Timestamp > end {
			continue
		}
		ret[res.Endpoint+res.Counter] = float64(v.Value)
	}

	return ret, nil
}
