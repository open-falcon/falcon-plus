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

type Token struct {
	Id         int64     `json:"id"`
	Name       string    `json:"name"`
	Cname      string    `json:"cname"`
	Note       string    `json:"note"`
	CreateTime time.Time `json:"ctime"`
}

func (op *Operator) AddToken(o *Token) (id int64, err error) {
	o.Id = 0
	if id, err = op.O.Insert(o); err != nil {
		return
	}
	o.Id = id
	moduleCache[CTL_M_TOKEN].set(id, o)
	DbLog(op.O, op.User.Id, CTL_M_TOKEN, id, CTL_A_ADD, jsonStr(o))
	return
}

func (op *Operator) GetToken(id int64) (o *Token, err error) {
	var ok bool

	if o, ok = moduleCache[CTL_M_TOKEN].get(id).(*Token); ok {
		return
	}
	o = &Token{Id: id}
	err = op.O.Read(o, "Id")
	if err == nil {
		moduleCache[CTL_M_TOKEN].set(id, o)
	}
	return
}

func (op *Operator) GetTokenByName(token string) (o *Token, err error) {
	o = &Token{Name: token}
	err = op.O.Read(o, "Name")
	return
}

func (op *Operator) QueryTokens(query string) orm.QuerySeter {
	qs := op.O.QueryTable(new(Token))
	if query != "" {
		qs = qs.Filter("Name__icontains", query)
	}
	return qs
}

func (op *Operator) GetTokensCnt(query string) (int64, error) {
	return op.QueryTokens(query).Count()
}

func (op *Operator) GetTokens(query string, limit, offset int) (tokens []*Token, err error) {
	_, err = op.QueryTokens(query).Limit(limit, offset).All(&tokens)
	return
}

func (op *Operator) UpdateToken(id int64, _tk *Token) (tk *Token, err error) {
	if tk, err = op.GetToken(id); err != nil {
		return nil, ErrNoToken
	}

	if _tk.Name != "" {
		tk.Name = _tk.Name
	}
	if _tk.Cname != "" {
		tk.Cname = _tk.Cname
	}
	if _tk.Note != "" {
		tk.Note = _tk.Note
	}
	_, err = op.O.Update(tk)
	moduleCache[CTL_M_TOKEN].set(id, tk)
	DbLog(op.O, op.User.Id, CTL_M_TOKEN, id, CTL_A_SET, "")
	return tk, err
}

func (op *Operator) DeleteToken(id int64) error {

	if n, err := op.O.Delete(&Token{Id: id}); err != nil || n == 0 {
		return ErrNoExits
	}
	moduleCache[CTL_M_TOKEN].del(id)
	DbLog(op.O, op.User.Id, CTL_M_TOKEN, id, CTL_A_DEL, "")

	return nil
}
