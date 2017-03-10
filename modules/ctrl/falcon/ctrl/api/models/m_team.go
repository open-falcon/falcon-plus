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
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type Team struct {
	Id         int64     `json:"id"`
	Name       string    `json:"name"`
	Note       string    `json:"note"`
	Creator    int64     `json:"-"`
	CreateTime time.Time `json:"ctime"`
}

type TeamMemberIds struct {
	Uids []int64 `json:"uids"`
}

type TeamMembers struct {
	Users []User `json:"users"`
}

func (op *Operator) AddTeam(t *Team) (id int64, err error) {
	t.Id = 0
	id, err = op.O.Insert(t)
	if err != nil {
		beego.Error(err)
		return
	}

	t.Id = id
	moduleCache[CTL_M_TEAM].set(id, t)
	DbLog(op.O, op.User.Id, CTL_M_TEAM, id, CTL_A_ADD, jsonStr(t))
	return
}

func (op *Operator) GetTeam(id int64) (*Team, error) {
	if t, ok := moduleCache[CTL_M_TEAM].get(id).(*Team); ok {
		return t, nil
	}
	t := &Team{Id: id}
	err := op.O.Read(t, "Id")
	if err == nil {
		moduleCache[CTL_M_TEAM].set(id, t)
	}
	return t, err
}

func (op *Operator) GetMember(id int64) (*TeamMembers, error) {
	var m TeamMembers

	_, err := op.O.Raw("SELECT `b`.`id`, `b`.`name` "+
		"FROM `team_user` `a` LEFT JOIN `user` `b` "+
		"ON `a`.`user_id` = `b`.`id` WHERE `a`.`team_id` = ? ",
		id).QueryRows(&m.Users)
	return &m, err
}

func (op *Operator) QueryTeams(query string, own bool) orm.QuerySeter {
	qs := op.O.QueryTable(new(Team))
	if query != "" {
		qs = qs.Filter("Name__icontains", query)
	}
	if own {
		qs = qs.Filter("Creator", op.User.Id)
	}
	return qs
}

func (op *Operator) GetTeamsCnt(query string, own bool) (int64, error) {
	return op.QueryTeams(query, own).Count()
}

func (op *Operator) GetTeams(query string, own bool, limit, offset int) (teams []*Team, err error) {
	_, err = op.QueryTeams(query, own).Limit(limit, offset).All(&teams)
	return
}

func (op *Operator) UpdateTeam(id int64, _t *Team) (t *Team, err error) {
	if t, err = op.GetTeam(id); err != nil {
		return nil, ErrNoTeam
	}

	if _t.Name != "" {
		t.Name = _t.Name
	}
	if _t.Note != "" {
		t.Note = _t.Note
	}
	_, err = op.O.Update(t)
	moduleCache[CTL_M_TEAM].set(id, t)
	DbLog(op.O, op.User.Id, CTL_M_TEAM, id, CTL_A_SET, jsonStr(_t))
	return t, err
}

func (op *Operator) UpdateMember(id int64, _m *TeamMemberIds) (m *TeamMemberIds, err error) {
	var tm *TeamMembers

	if tm, err = op.GetMember(id); err != nil {
		return nil, ErrNoTeam
	}

	m = &TeamMemberIds{}
	m.Uids = make([]int64, len(tm.Users))
	for i, v := range tm.Users {
		m.Uids[i] = v.Id
	}

	add, del := MdiffInt(m.Uids, _m.Uids)
	if len(add) > 0 {
		vs := make([]string, len(add))
		for i := 0; i < len(vs); i++ {
			vs[i] = fmt.Sprintf("(%d, %d)", id, add[i])
		}
		if _, err = op.O.Raw("INSERT `team_user` (`team_id`, `user_id`) VALUES " + strings.Join(vs, ", ")).Exec(); err != nil {
			return
		}
	}
	if len(del) > 0 {
		ids := fmt.Sprintf("%d", del[0])
		for i := 0; i < len(del)-1; i++ {
			ids += fmt.Sprintf("%s, %d", ids, del[i])
		}
		if _, err = op.O.Raw("DELETE from `team_user` WHERE team_id = ? and user_id in ("+ids+")", id).Exec(); err != nil {
			return
		}
	}
	m.Uids = _m.Uids
	DbLog(op.O, op.User.Id, CTL_M_TEAM, id, CTL_A_SET, jsonStr(_m))
	return
}

func (op *Operator) DeleteTeam(id int64) error {
	if n, err := op.O.Delete(&Team{Id: id}); err != nil || n == 0 {
		return err
	}
	DbLog(op.O, op.User.Id, CTL_M_TEAM, id, CTL_A_DEL, "")

	return nil
}

func (op *Operator) BindTeamUser(team_id, user_id int64) (err error) {
	if _, err := op.O.Raw("INSERT INTO `team_user` (`team_id`, `user_id`) VALUES (?, ?)", team_id, user_id).Exec(); err != nil {
		return err
	}
	return nil
}

func (op *Operator) BindTeamUsers(team_id int64, user_ids []int64) (int64, error) {
	vs := make([]string, len(user_ids))
	for i := 0; i < len(vs); i++ {
		vs[i] = fmt.Sprintf("(%d, %d)", team_id, user_ids[i])
	}

	if res, err := op.O.Raw("INSERT `team_user` (`team_id`, `user_id`) VALUES " + strings.Join(vs, ", ")).Exec(); err != nil {
		return 0, err
	} else {
		return res.RowsAffected()
	}
}
