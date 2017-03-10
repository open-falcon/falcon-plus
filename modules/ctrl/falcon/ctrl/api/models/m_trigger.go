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

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type Trigger struct {
	Id         int64     `json:"id"`
	Name       string    `json:"name"`
	Cname      string    `json:"cname"`
	Note       string    `json:"note"`
	CreateTime time.Time `json:"ctime"`
}

func (op *Operator) AddTrigger(r *Trigger) (id int64, err error) {
	r.Id = 0
	id, err = op.O.Insert(r)
	if err != nil {
		return
	}
	r.Id = id
	moduleCache[CTL_M_TRIGGER].set(id, r)
	DbLog(op.O, op.User.Id, CTL_M_TRIGGER, id, CTL_A_ADD, jsonStr(r))
	return
}

func (op *Operator) GetTrigger(id int64) (*Trigger, error) {
	if r, ok := moduleCache[CTL_M_TRIGGER].get(id).(*Trigger); ok {
		return r, nil
	}
	r := &Trigger{Id: id}
	err := op.O.Read(r, "Id")
	if err == nil {
		moduleCache[CTL_M_TRIGGER].set(id, r)
	}
	return r, err
}

func (op *Operator) QueryTriggers(query string) orm.QuerySeter {
	qs := op.O.QueryTable(new(Trigger))
	if query != "" {
		qs = qs.Filter("Name__icontains", query)
	}
	return qs
}

func (op *Operator) GetTriggersCnt(query string) (int64, error) {
	return op.QueryTriggers(query).Count()
}

func (op *Operator) GetTriggers(query string, limit, offset int) (triggers []*Trigger, err error) {
	_, err = op.QueryTriggers(query).Limit(limit, offset).All(&triggers)
	return
}

func (op *Operator) UpdateTrigger(id int64, _r *Trigger) (r *Trigger, err error) {
	if r, err = op.GetTrigger(id); err != nil {
		return nil, ErrNoTrigger
	}

	if _r.Name != "" {
		r.Name = _r.Name
	}
	if _r.Cname != "" {
		r.Cname = _r.Cname
	}
	if _r.Note != "" {
		r.Note = _r.Note
	}
	_, err = op.O.Update(r)
	moduleCache[CTL_M_TRIGGER].set(id, r)
	DbLog(op.O, op.User.Id, CTL_M_TRIGGER, id, CTL_A_SET, "")
	return r, err
}

func (op *Operator) DeleteTrigger(id int64) error {
	if n, err := op.O.Delete(&Trigger{Id: id}); err != nil || n == 0 {
		return err
	}
	moduleCache[CTL_M_TRIGGER].del(id)
	DbLog(op.O, op.User.Id, CTL_M_TRIGGER, id, CTL_A_DEL, "")

	return nil
}

func (op *Operator) BindUserTrigger(user_id, trigger_id, tag_id int64) (err error) {
	if _, err := op.O.Raw("INSERT INTO `tag_trigger_user` (`tag_id`, `trigger_id`, `user_id`) VALUES (?, ?, ?)", tag_id, trigger_id, user_id).Exec(); err != nil {
		return err
	}
	return nil
}

func (op *Operator) BindTokenTrigger(token_id, trigger_id, tag_id int64) (err error) {
	if _, err := op.O.Raw("INSERT INTO `tag_trigger_token` (`tag_id`, `trigger_id`, `token_id`) VALUES (?, ?, ?)", tag_id, trigger_id, token_id).Exec(); err != nil {
		return err
	}
	return nil
}
