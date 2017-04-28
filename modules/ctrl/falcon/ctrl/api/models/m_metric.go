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
package models

import "strings"

type Metric struct {
	Name string `json:"name"`
}

var (
	metrics []*Metric
)

func queryMetrics(query string, limit, offset int) (ret []*Metric) {
	for k, v := range metrics {
		if strings.Contains(v.Name, query) {
			if offset == 0 {
				ret = append(ret, metrics[k])
			} else {
				offset--
			}
			if limit == 0 {
				return
			} else {
				limit--
			}
		}
	}
	return
}

func (op *Operator) GetMetricsCnt(query string) (int64, error) {
	return int64(len(queryMetrics(query, -1, 0))), nil
}

func (op *Operator) GetMetrics(query string, limit, offset int) (metrics []*Metric, err error) {
	return queryMetrics(query, limit, offset), nil
}
