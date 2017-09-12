// Copyright 2017 Xiaomi, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
