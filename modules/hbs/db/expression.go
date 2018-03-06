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
	"log"
	"strings"
)

func QueryExpressions() (ret []*model.Expression, err error) {
	sql := "select * from expression where action_id>0 and pause=0"
	rows, err := DB.Query(sql)
	if err != nil {
		log.Println("ERROR:", err)
		return ret, err
	}

	defer rows.Close()
	for rows.Next() {
		e := model.Expression{}
		var exp string
		err = rows.Scan(
			&e.Id,
			&exp,
			&e.Func,
			&e.Operator,
			&e.RightValue,
			&e.MaxStep,
			&e.Priority,
			&e.Note,
			&e.ActionId,
			&e.CreateUser,
			&e.Pause,
		)

		if err != nil {
			log.Println("WARN:", err)
			continue
		}

		e.Metric, e.Tags, err = parseExpression(exp)
		if err != nil {
			log.Println("ERROR:", err)
			continue
		}

		ret = append(ret, &e)
	}

	return ret, nil
}

func parseExpression(exp string) (metric string, tags map[string]string, err error) {
	left := strings.Index(exp, "(")
	right := strings.Index(exp, ")")
	tagStrs := strings.TrimSpace(exp[left+1 : right])

	arr := strings.Fields(tagStrs)
	if len(arr) < 2 {
		err = fmt.Errorf("tag not enough. exp: %s", exp)
		return
	}

	tags = make(map[string]string)
	for _, item := range arr {
		kv := strings.Split(item, "=")
		if len(kv) != 2 {
			err = fmt.Errorf("parse %s fail", exp)
			return
		}
		tags[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
	}

	metric, exists := tags["metric"]
	if !exists {
		err = fmt.Errorf("no metric give of %s", exp)
		return
	}

	delete(tags, "metric")
	return
}
