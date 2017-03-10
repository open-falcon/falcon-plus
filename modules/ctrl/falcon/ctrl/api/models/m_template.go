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

import "time"

type TemplateAction struct {
	Template Template `json:"template"`
	Action   Action   `json:"action"`
	Pname    string   `json:"pname"`
}

// for ui
type TemplateUi struct {
	Id      int64  `json:"id"`
	Pid     int64  `json:"pid"`
	Name    string `json:"name"`
	Pname   string `json:"pname"`
	Creator string `json:"creator"`
}

// for db
type Template struct {
	Id           int64     `json:"id"`
	Name         string    `json:"name"`
	ParentId     int64     `json:"pid"`
	ActionId     int64     `json:"-"`
	CreateUserId int64     `json:"-"`
	CreateTime   time.Time `json:"ctime"`
}

// for db
type Action struct {
	Id                 int64  `json:"id"`
	Uic                string `json:"uic"`
	Url                string `json:"url"`
	SendSms            uint   `json:"sendSms"`
	SendMail           uint   `json:"sendMail"`
	Callback           uint   `json:"callback"`
	BeforeCallbackSms  uint   `json:"beforeCallbackSms"`
	BeforeCallbackMail uint   `json:"beforeCallbackMail"`
	AfterCallbackSms   uint   `json:"afterCallbackSms"`
	AfterCallbackMail  uint   `json:"afterCallbackMail"`
}

func (op *Operator) AddTemplate(o *Template) (id int64, err error) {
	o.CreateUserId = op.User.Id
	o.Id = 0
	id, err = op.O.Insert(o)
	if err != nil {
		return
	}
	o.Id = id
	moduleCache[CTL_M_TEMPLATE].set(id, o)
	DbLog(op.O, op.User.Id, CTL_M_TEMPLATE, id, CTL_A_ADD, jsonStr(o))
	return
}

func (op *Operator) AddAction(o *Action) (id int64, err error) {
	o.Id = 0
	id, err = op.O.Insert(o)
	if err != nil {
		return
	}
	o.Id = id
	return
}

func (op *Operator) CloneTemplate(id int64) (*TemplateAction, error) {
	var (
		ret     TemplateAction
		src_tpl *Template
		src_act *Action
		err     error
		objs    []*Strategy
		tid     int64
	)

	if src_tpl, err = op.getTemplate(id); err != nil {
		return nil, err
	}
	if src_act, err = op.GetAction(src_tpl.ActionId); err != nil {
		return nil, err
	}

	ret.Template = *src_tpl
	ret.Action = *src_act
	ret.Template.Name += "_copy"

	_, err = op.AddAction(&ret.Action)
	if err != nil {
		return nil, err
	}

	ret.Template.ActionId = ret.Action.Id
	tid, err = op.AddTemplate(&ret.Template)
	if err != nil {
		return nil, err
	}

	objs, _ = op.GetStrategys(src_tpl.Id, "", 0, 0)
	for _, obj := range objs {
		obj.TplId = tid
		op.AddStrategy(obj)
	}

	if t, err := op.getTemplate(ret.Template.ParentId); err == nil {
		ret.Pname = t.Name
	}

	return &ret, nil
}

func (op *Operator) GetTemplate(id int64) (*TemplateAction, error) {
	var ret TemplateAction

	if t, err := op.getTemplate(id); err != nil {
		return nil, err
	} else {
		ret.Template = *t
	}

	if a, err := op.GetAction(ret.Template.ActionId); err != nil {
		return nil, err
	} else {
		ret.Action = *a
	}

	if t, err := op.getTemplate(ret.Template.ParentId); err == nil {
		ret.Pname = t.Name
	}

	return &ret, nil
}

func (op *Operator) getTemplate(id int64) (*Template, error) {
	if r, ok := moduleCache[CTL_M_TEMPLATE].get(id).(*Template); ok {
		return r, nil
	}
	r := &Template{Id: id}
	err := op.O.Read(r, "Id")
	if err == nil {
		moduleCache[CTL_M_TEMPLATE].set(id, r)
	}
	return r, err
}

func (op *Operator) GetAction(id int64) (*Action, error) {
	a := &Action{Id: id}
	err := op.O.Read(a, "Id")
	return a, err
}

func (op *Operator) GetTemplatesCnt(query string) (cnt int64, err error) {
	if query == "" {
		err = op.O.Raw("SELECT count(*) FROM template").QueryRow(&cnt)
	} else {
		err = op.O.Raw("SELECT count(*) FROM template WHERE name like ?", "%"+query+"%").QueryRow(&cnt)
	}
	return
}

func (op *Operator) GetTemplates(query string, limit, offset int) (ret []TemplateUi, err error) {
	if query == "" {
		_, err = op.O.Raw("SELECT a.id, b.id as pid, a.name, b.name as pname, c.name as creator  FROM template a LEFT JOIN template b ON a.parent_id = b.id LEFT JOIN user c ON a.create_user_id = c.id ORDER BY a.name LIMIT ? OFFSET ?", limit, offset).QueryRows(&ret)
	} else {
		_, err = op.O.Raw("SELECT a.id as id, b.id as pid, a.name as name, b.name as pname, c.name as creator  FROM template a LEFT JOIN template b ON a.parent_id = b.id LEFT JOIN user c ON a.create_user_id = c.id WHERE a.name like ? ORDER BY a.name LIMIT ? OFFSET ?", "%"+query+"%", limit, offset).QueryRows(&ret)
	}
	return
}

func (op *Operator) UpdateTemplate(id int64, _o *TemplateAction) (o *Template, err error) {
	var t *Template
	t, err = op.updateTemplate(id, &_o.Template)
	if err != nil {
		return nil, err
	}
	_, err = op.UpdateAction(t.ActionId, &_o.Action)
	if err != nil {
		return nil, err
	}
	return t, err
}

func (op *Operator) updateTemplate(id int64, _o *Template) (o *Template, err error) {
	if o, err = op.getTemplate(id); err != nil {
		return nil, ErrNoTemplate
	}

	if _o.Name != "" {
		o.Name = _o.Name
	}
	if _o.ParentId != 0 {
		o.ParentId = _o.ParentId
	}
	_, err = op.O.Update(o)
	return o, err
}

func (op *Operator) UpdateAction(id int64, _o *Action) (o *Action, err error) {
	if o, err = op.GetAction(id); err != nil {
		return nil, ErrNoTemplate
	}
	if _o.Uic != "" {
		o.Uic = _o.Uic
	}
	if _o.Url != "" {
		o.Url = _o.Url
	}
	o.SendSms = _o.SendSms
	o.SendMail = _o.SendMail
	o.Callback = _o.Callback
	o.BeforeCallbackSms = _o.BeforeCallbackSms
	o.BeforeCallbackMail = _o.BeforeCallbackMail
	o.AfterCallbackSms = _o.AfterCallbackSms
	o.AfterCallbackMail = _o.AfterCallbackMail
	_, err = op.O.Update(o)
	return o, err
}

func (op *Operator) DeleteTemplate(id int64) error {
	template, err := op.getTemplate(id)
	if err != nil {
		return err
	}

	if _, err = op.O.Delete(&Action{Id: template.ActionId}); err != nil {
		return err
	}

	if _, err = op.O.Delete(&Template{Id: id}); err != nil {
		return err
	}

	moduleCache[CTL_M_TEMPLATE].del(id)
	DbLog(op.O, op.User.Id, CTL_M_TEMPLATE, id, CTL_A_DEL, "")

	return nil
}
