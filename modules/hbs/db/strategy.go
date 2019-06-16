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

package db

import (
	"fmt"
	"github.com/open-falcon/falcon-plus/common/model"
	"github.com/toolkits/container/set"
	"log"
	"strings"
	"time"
)

// 获取所有的Strategy列表
func QueryStrategies(tpls map[int]*model.Template) (map[int]*model.Strategy, error) {
	ret := make(map[int]*model.Strategy)

	if tpls == nil || len(tpls) == 0 {
		return ret, fmt.Errorf("illegal argument")
	}

	now := time.Now().Format("15:04")
	sql := fmt.Sprintf(
		"select %s from strategy as s where (s.run_begin='' and s.run_end='') "+
			"or (s.run_begin <= '%s' and s.run_end >= '%s')"+
			"or (s.run_begin > s.run_end and !(s.run_begin > '%s' and s.run_end < '%s'))",
		"s.id, s.metric, s.tags, s.func, s.op, s.right_value, s.max_step, s.priority, s.note, s.tpl_id",
		now,
		now,
		now,
		now,
	)

	rows, err := DB.Query(sql)
	if err != nil {
		log.Println("ERROR:", err)
		return ret, err
	}

	defer rows.Close()
	for rows.Next() {
		s := model.Strategy{}
		var tags string
		var tid int
		err = rows.Scan(&s.Id, &s.Metric, &tags, &s.Func, &s.Operator, &s.RightValue, &s.MaxStep, &s.Priority, &s.Note, &tid)
		if err != nil {
			log.Println("ERROR:", err)
			continue
		}

		tt := make(map[string]string)

		if tags != "" {
			arr := strings.Split(tags, ",")
			for _, tag := range arr {
				kv := strings.SplitN(tag, "=", 2)
				if len(kv) != 2 {
					continue
				}
				tt[kv[0]] = kv[1]
			}
		}

		s.Tags = tt
		s.Tpl = tpls[tid]
		if s.Tpl == nil {
			log.Printf("WARN: tpl is nil. strategy id=%d, tpl id=%d", s.Id, tid)
			// 如果Strategy没有对应的Tpl，那就没有action，就没法报警，无需往后传递了
			continue
		}

		ret[s.Id] = &s
	}

	return ret, nil
}

func QueryBuiltinMetrics(tids string) ([]*model.BuiltinMetric, error) {
	sql := fmt.Sprintf(
		"select metric, tags from strategy where tpl_id in (%s) and metric in ('net.port.listen', 'proc.num', 'du.bs', 'url.check.health')",
		tids,
	)

	ret := []*model.BuiltinMetric{}

	rows, err := DB.Query(sql)
	if err != nil {
		log.Println("ERROR:", err)
		return ret, err
	}

	metricTagsSet := set.NewStringSet()

	defer rows.Close()
	for rows.Next() {
		builtinMetric := model.BuiltinMetric{}
		err = rows.Scan(&builtinMetric.Metric, &builtinMetric.Tags)
		if err != nil {
			log.Println("WARN:", err)
			continue
		}

		k := fmt.Sprintf("%s%s", builtinMetric.Metric, builtinMetric.Tags)
		if metricTagsSet.Exists(k) {
			continue
		}

		ret = append(ret, &builtinMetric)
		metricTagsSet.Add(k)
	}

	return ret, nil
}
