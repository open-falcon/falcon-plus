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

import "github.com/astaxie/beego/orm"

type Strategy struct {
	Id        int64  `json:"id"`
	MetricId  int64  `json:"metricId"`
	Tags      string `json:"tags"`
	MaxStep   int    `json:"maxStep"`
	Priority  int    `json:"priority"`
	Func      string `json:"fun"`
	Op        string `json:"op"`
	Condition string `json:"condition"`
	Note      string `json:"note"`
	Metric    string `json:"metric"`
	RunBegin  string `json:"runBegin"`
	RunEnd    string `json:"runEnd"`
	TplId     int64  `json:"tplId"`
}

func (op *Operator) AddStrategy(o *Strategy) (id int64, err error) {
	o.Id = 0
	if id, err = op.O.Insert(o); err != nil {
		return
	}
	o.Id = id
	return
}

func (op *Operator) GetStrategy(id int64) (s *Strategy, err error) {
	s = &Strategy{Id: id}
	err = op.O.Read(s, "Id")
	return
}

func (op *Operator) QueryStrategys(tid int64, query string) orm.QuerySeter {
	qs := op.O.QueryTable(new(Strategy))
	if tid != 0 {
		qs = qs.Filter("TplId", tid)
	}
	if query != "" {
		qs = qs.Filter("Name__icontains", query)
	}
	return qs
}

func (op *Operator) GetStrategysCnt(tid int64, query string) (int64, error) {
	return op.QueryStrategys(tid, query).Count()
}

func (op *Operator) GetStrategys(tid int64, query string, limit, offset int) (strategys []*Strategy, err error) {
	_, err = op.QueryStrategys(tid, query).Limit(limit, offset).All(&strategys)
	return
}

func (op *Operator) UpdateStrategy(id int64, _o *Strategy) (o *Strategy, err error) {
	if o, err = op.GetStrategy(id); err != nil {
		return nil, ErrNoStrategy
	}

	o.MetricId = _o.MetricId
	o.Tags = _o.Tags
	o.MaxStep = _o.MaxStep
	o.Priority = _o.Priority
	o.Func = _o.Func
	o.Op = _o.Op
	o.Condition = _o.Condition
	o.Note = _o.Note
	o.Metric = _o.Metric
	o.RunBegin = _o.RunBegin
	o.RunEnd = _o.RunEnd
	o.TplId = _o.TplId

	_, err = op.O.Update(o)
	return o, err
}

func (op *Operator) DeleteStrategy(id int64) error {

	if n, err := op.O.Delete(&Strategy{Id: id}); err != nil || n == 0 {
		return ErrNoExits
	}

	return nil
}
