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

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type Host struct {
	Id         int64     `json:"id"`
	Uuid       string    `json:"uuid"`
	Name       string    `json:"name"`
	Type       string    `json:"typ"`
	Status     string    `json:"status"`
	Loc        string    `json:"loc"`
	Idc        string    `json:"idc"`
	CreateTime time.Time `json:"ctime"`
}

func (op *Operator) AddHost(h *Host) (id int64, err error) {
	h.Id = 0
	id, err = op.O.Insert(h)
	if err != nil {
		beego.Error(err)
		return
	}
	h.Id = id
	moduleCache[CTL_M_HOST].set(id, h)
	DbLog(op.O, op.User.Id, CTL_M_HOST, id, CTL_A_ADD, jsonStr(h))
	return
}

func (op *Operator) GetHost(id int64) (*Host, error) {
	if h, ok := moduleCache[CTL_M_HOST].get(id).(*Host); ok {
		return h, nil
	}
	h := &Host{Id: id}
	err := op.O.Read(h, "Id")
	if err == nil {
		moduleCache[CTL_M_HOST].set(id, h)
	}
	return h, err
}

func (op *Operator) GetHostByUuid(uuid string) (h *Host, err error) {
	h = &Host{Uuid: uuid}
	err = op.O.Read(h, "Uuid")
	return h, err
}

func (op *Operator) QueryHosts(query string) orm.QuerySeter {
	// TODO: acl filter
	// just for admin?
	qs := op.O.QueryTable(new(Host))
	if query != "" {
		qs = qs.Filter("Name__icontains", query)
	}
	return qs
}

func (op *Operator) GetHostsCnt(query string) (int64, error) {
	return op.QueryHosts(query).Count()
}

func (op *Operator) GetHosts(query string, limit, offset int) (hosts []*Host, err error) {
	_, err = op.QueryHosts(query).Limit(limit, offset).All(&hosts)
	return
}

func (op *Operator) UpdateHost(id int64, _h *Host) (h *Host, err error) {
	if h, err = op.GetHost(id); err != nil {
		return nil, ErrNoHost
	}

	if _h.Uuid != "" {
		h.Uuid = _h.Uuid
	}
	if _h.Name != "" {
		h.Name = _h.Name
	}
	if _h.Type != "" {
		h.Type = _h.Type
	}
	if _h.Type != "" {
		h.Type = _h.Type
	}
	if _h.Status != "" {
		h.Status = _h.Status
	}
	if _h.Loc != "" {
		h.Loc = _h.Loc
	}
	if _h.Idc != "" {
		h.Idc = _h.Idc
	}
	_, err = op.O.Update(h)
	moduleCache[CTL_M_HOST].set(id, h)
	DbLog(op.O, op.User.Id, CTL_M_HOST, id, CTL_A_SET, "")
	return h, err
}

func (op *Operator) DeleteHost(id int64) error {
	if n, err := op.O.Delete(&Host{Id: id}); err != nil || n == 0 {
		return err
	}
	moduleCache[CTL_M_HOST].del(id)
	DbLog(op.O, op.User.Id, CTL_M_HOST, id, CTL_A_DEL, "")

	return nil
}
