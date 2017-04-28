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

type Role struct {
	Id         int64     `json:"id"`
	Name       string    `json:"name"`
	Cname      string    `json:"cname"`
	Note       string    `json:"note"`
	CreateTime time.Time `json:"ctime"`
}

func (op *Operator) AddRole(r *Role) (id int64, err error) {
	r.Id = 0
	id, err = op.O.Insert(r)
	if err != nil {
		return
	}
	r.Id = id
	moduleCache[CTL_M_ROLE].set(id, r)
	DbLog(op.O, op.User.Id, CTL_M_ROLE, id, CTL_A_ADD, jsonStr(r))
	return
}

func (op *Operator) GetRole(id int64) (*Role, error) {
	if r, ok := moduleCache[CTL_M_ROLE].get(id).(*Role); ok {
		return r, nil
	}
	r := &Role{Id: id}
	err := op.O.Read(r, "Id")
	if err == nil {
		moduleCache[CTL_M_ROLE].set(id, r)
	}
	return r, err
}

func (op *Operator) QueryRoles(query string) orm.QuerySeter {
	qs := op.O.QueryTable(new(Role))
	if query != "" {
		qs = qs.Filter("Name__icontains", query)
	}
	return qs
}

func (op *Operator) GetRolesCnt(query string) (int64, error) {
	return op.QueryRoles(query).Count()
}

func (op *Operator) GetRoles(query string, limit, offset int) (roles []*Role, err error) {
	_, err = op.QueryRoles(query).Limit(limit, offset).All(&roles)
	return
}

func (op *Operator) UpdateRole(id int64, _r *Role) (r *Role, err error) {
	if r, err = op.GetRole(id); err != nil {
		return nil, ErrNoRole
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
	moduleCache[CTL_M_ROLE].set(id, r)
	DbLog(op.O, op.User.Id, CTL_M_ROLE, id, CTL_A_SET, "")
	return r, err
}

func (op *Operator) DeleteRole(id int64) error {
	if n, err := op.O.Delete(&Role{Id: id}); err != nil || n == 0 {
		return err
	}
	moduleCache[CTL_M_ROLE].del(id)
	DbLog(op.O, op.User.Id, CTL_M_ROLE, id, CTL_A_DEL, "")

	return nil
}

func (op *Operator) BindUserRole(user_id, role_id, tag_id int64) (err error) {
	if _, err := op.O.Raw("INSERT INTO `tag_role_user` (`tag_id`, `role_id`, `user_id`) VALUES (?, ?, ?)", tag_id, role_id, user_id).Exec(); err != nil {
		return err
	}
	return nil
}

func (op *Operator) BindTokenRole(token_id, role_id, tag_id int64) (err error) {
	if _, err := op.O.Raw("INSERT INTO `tag_role_token` (`tag_id`, `role_id`, `token_id`) VALUES (?, ?, ?)", tag_id, role_id, token_id).Exec(); err != nil {
		return err
	}
	return nil
}
