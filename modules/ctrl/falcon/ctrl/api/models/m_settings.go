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

	"github.com/open-falcon/falcon-plus/modules/ctrl/falcon/ctrl"
)

type LogUi struct {
	LogId    int64     `json:"log_id"`
	Module   int64     `json:"module"`
	Id       int64     `json:"id"`
	User     string    `json:"user"`
	ActionId int64     `json:"action_id"`
	Data     string    `json:"data"`
	Time     time.Time `json:"time"`
}

func logSql(begin, end string) (where string, args []interface{}) {
	sql2 := []string{}
	sql3 := []interface{}{}
	if begin != "" {
		sql2 = append(sql2, "a.time >= ?")
		sql3 = append(sql3, begin)
	}
	if end != "" {
		sql2 = append(sql2, "a.time <= ?")
		sql3 = append(sql3, end)
	}
	if len(sql2) != 0 {
		where = "WHERE " + strings.Join(sql2, " AND ")
		args = sql3
	}
	return
}

func (op *Operator) GetLogsCnt(begin, end string) (cnt int64, err error) {
	sql, sql_args := logSql(begin, end)
	err = op.O.Raw("SELECT count(*) FROM log a "+sql, sql_args...).QueryRow(&cnt)
	return
}

func (op *Operator) GetLogs(begin, end string, limit, offset int) (ret []*LogUi, err error) {
	sql, sql_args := logSql(begin, end)
	sql = "select a.id as log_id, a.module, a.module_id as id, b.name as user, a.action, a.data, a.time from log a left join user b on a.user_id = b.id " + sql + " ORDER BY a.id DESC LIMIT ? OFFSET ?"
	sql_args = append(sql_args, limit, offset)
	_, err = op.O.Raw(sql, sql_args...).QueryRows(&ret)
	return
}

func (op *Operator) Populate() (interface{}, error) {
	var (
		ret       string
		err       error
		items     []string
		user      *User
		id        int64
		tag_idx   = make(map[string]int64)
		user_idx  = make(map[string]int64)
		team_idx  = make(map[string]int64)
		role_idx  = make(map[string]int64)
		token_idx = make(map[string]int64)
		host_idx  = make(map[string]int64)
		tpl_idx   = make(map[string]int64)
		test_user = "test01"
	)
	tag_idx["/"] = 1

	// user
	items = []string{
		test_user,
		"user0",
		"user1",
		"user2",
		"user3",
		"user4",
		"user5",
		"user6",
	}
	for _, item := range items {
		if user, err = op.AddUser(&User{Name: item, Uuid: item}); err != nil {
			return nil, err
		}
		user_idx[item] = user.Id
		ret = fmt.Sprintf("%sadd user(%s)\n", ret, item)
	}

	// team
	items = []string{
		"team1",
		"team2",
		"team3",
		"team4",
	}
	for _, item := range items {
		if id, err = op.AddTeam(&Team{Name: item, Creator: op.User.Id}); err != nil {
			fmt.Printf("add team(%s)\n", item)
			return nil, err
		}
		team_idx[item] = id
		ret = fmt.Sprintf("%sadd team(%s)\n", ret, item)
	}
	teamMembers := []struct {
		team  string
		users []string
	}{
		{"team1", []string{"user0", "user1"}},
		{"team2", []string{"user2", "user3"}},
		{"team3", []string{"user4", "user5"}},
		{"team4", []string{"user0", "user1", "user2", "user3", "user4", "user5"}},
	}
	for _, item := range teamMembers {
		uids := make([]int64, len(item.users))
		for i := 0; i < len(uids); i++ {
			uids[i] = user_idx[item.users[i]]
		}
		if _, err = op.UpdateMember(team_idx[item.team],
			&TeamMemberIds{Uids: uids}); err != nil {
			return nil, err
		}
		ret = fmt.Sprintf("%sadd teamMembers(%v)\n", ret, item)
	}

	// tag
	items = []string{
		"cop=xiaomi",
		"cop=xiaomi,owt=inf",
		"cop=xiaomi,owt=miliao",
		"cop=xiaomi,owt=miliao,pdl=op",
		"cop=xiaomi,owt=miliao,pdl=micloud",
	}
	for _, item := range items {
		if tag_idx[item], err = op.AddTag(&Tag{Name: item}); err != nil {
			return nil, err
		}
		ret = fmt.Sprintf("%sadd tag(%s)\n", ret, item)
	}

	// tag host
	items2 := [][2]string{
		{"cop=xiaomi", "mi1.bj"},
		{"cop=xiaomi", "mi2.bj"},
		{"cop=xiaomi", "mi3.bj"},
		{"cop=xiaomi,owt=inf", "inf1.bj"},
		{"cop=xiaomi,owt=inf", "inf2.bj"},
		{"cop=xiaomi,owt=inf", "inf3.bj"},
		{"cop=xiaomi,owt=miliao", "miliao1.bj"},
		{"cop=xiaomi,owt=miliao", "miliao2.bj"},
		{"cop=xiaomi,owt=miliao", "miliao3.bj"},
		{"cop=xiaomi,owt=miliao,pdl=op", "miliao.op1.bj"},
		{"cop=xiaomi,owt=miliao,pdl=op", "miliao.op2.bj"},
		{"cop=xiaomi,owt=miliao,pdl=op", "miliao.op3.bj"},
		{"cop=xiaomi,owt=miliao,pdl=micloud", "miliao.cloud1.bj"},
		{"cop=xiaomi,owt=miliao,pdl=micloud", "miliao.cloud2.bj"},
		{"cop=xiaomi,owt=miliao,pdl=micloud", "miliao.cloud3.bj"},
	}
	for _, item2 := range items2 {
		if host_idx[item2[1]], err = op.AddHost(&Host{Name: item2[1]}); err != nil {
			return nil, err
		}

		if _, err = op.CreateTagHost(RelTagHost{TagId: tag_idx[item2[0]],
			HostId: host_idx[item2[1]]}); err != nil {
			return nil, err
		}
		ret = fmt.Sprintf("%sadd host(%s, %s)\n", ret, item2[1], item2[0])
	}

	// template
	items = []string{
		"tpl1",
		"tpl2",
		"tpl3",
	}
	for _, item := range items {
		if id, err = op.AddAction(&Action{}); err != nil {
			return nil, err
		}
		if tpl_idx[item], err = op.AddTemplate(&Template{Name: item,
			ActionId: id}); err != nil {
			return nil, err
		}
		ret = fmt.Sprintf("%sadd tag(%s)\n", ret, item)
	}
	// template strategy
	items2 = [][2]string{
		{"tpl1", "cpu.busy"},
		{"tpl1", "cpu.cnt"},
		{"tpl1", "cpu.idle"},
		{"tpl2", "cpu.busy"},
		{"tpl2", "cpu.cnt"},
		{"tpl2", "cpu.idle"},
		{"tpl3", "cpu.busy"},
		{"tpl3", "cpu.cnt"},
		{"tpl3", "cpu.idle"},
	}
	for _, item2 := range items2 {
		if _, err = op.AddStrategy(&Strategy{Metric: item2[1],
			TplId: tpl_idx[item2[0]]}); err != nil {
			return nil, err
		}
		ret = fmt.Sprintf("%sadd strategy(%s, %s)\n",
			ret, item2[0], item2[1])
	}

	// clone template
	items = []string{
		"tpl1",
		"tpl2",
		"tpl3",
	}
	for _, item := range items {
		if _, err = op.CloneTemplate(tpl_idx[item]); err != nil {
			return nil, err
		}
		ret = fmt.Sprintf("%s clone template(%s)\n", ret, item)
	}

	// bind tag template
	items2 = [][2]string{
		{"cop=xiaomi", "tpl1"},
		{"cop=xiaomi", "tpl2"},
		{"cop=xiaomi", "tpl3"},
		{"cop=xiaomi,owt=inf", "tpl1"},
		{"cop=xiaomi,owt=inf", "tpl2"},
		{"cop=xiaomi,owt=inf", "tpl3"},
		{"cop=xiaomi,owt=miliao,pdl=op", "tpl1"},
		{"cop=xiaomi,owt=miliao,pdl=op", "tpl2"},
		{"cop=xiaomi,owt=miliao,pdl=op", "tpl3"},
	}
	for _, item2 := range items2 {
		if _, err = op.CreateTagTpl(RelTagTpl{TagId: tag_idx[item2[0]],
			TplId: tpl_idx[item2[1]]}); err != nil {
			return nil, err
		}
		ret = fmt.Sprintf("%sadd tag tpl(%s, %s)\n", ret, item2[0], item2[1])
	}

	// role
	items = []string{
		"adm",
		"sre",
		"dev",
		"usr",
	}
	for _, item := range items {
		if role_idx[item], err = op.AddRole(&Role{Name: item}); err != nil {
			return nil, err
		}
		ret = fmt.Sprintf("%sadd role(%s)\n", ret, item)
	}

	// token
	items = []string{
		SYS_R_TOKEN,
		SYS_O_TOKEN,
		SYS_A_TOKEN,
	}
	for _, item := range items {
		if token_idx[item], err = op.AddToken(&Token{Name: item}); err != nil {
			return nil, err
		}
		ret = fmt.Sprintf("%sadd token(%s)\n", ret, item)
	}

	// bind user
	binds := [][3]string{
		{"cop=xiaomi,owt=miliao", test_user, "adm"},
		{"cop=xiaomi,owt=miliao", test_user, "sre"},
		{"cop=xiaomi,owt=miliao", test_user, "dev"},
		{"cop=xiaomi,owt=miliao", test_user, "usr"},
	}
	for _, s := range binds {
		if _, err := addTplRel(op.O, op.User.Id, tag_idx[s[0]], role_idx[s[2]],
			user_idx[s[1]], TPL_REL_T_ACL_USER); err != nil {
			return nil, err
		}
		ret = fmt.Sprintf("%sbind tag(%s) user(%s) role(%s)\n",
			ret, s[0], s[1], s[2])
	}

	// bind token
	binds = [][3]string{
		{SYS_O_TOKEN, "adm", "/"},
		{SYS_R_TOKEN, "adm", "/"},
		{SYS_A_TOKEN, "adm", "/"},
		{SYS_O_TOKEN, "sre", "/"},
		{SYS_R_TOKEN, "sre", "/"},
		{SYS_R_TOKEN, "dev", "/"},
		{SYS_R_TOKEN, "usr", "/"},
		{SYS_O_TOKEN, "adm", "cop=xiaomi,owt=miliao"},
		{SYS_O_TOKEN, "dev", "cop=xiaomi,owt=miliao,pdl=op"},
		{SYS_O_TOKEN, "usr", "cop=xiaomi"},
		{SYS_A_TOKEN, "usr", "cop=xiaomi,owt=miliao"},
	}
	for _, s := range binds {
		if _, err := addTplRel(op.O, op.User.Id, tag_idx[s[2]], role_idx[s[1]],
			token_idx[s[0]], TPL_REL_T_ACL_TOKEN); err != nil {
			return nil, err
		}
		ret = fmt.Sprintf("%sbind tag(%s) token(%s) role(%s)\n",
			ret, s[1], s[2], s[0])
	}

	return ret, nil
}

func (op *Operator) ResetDb() (interface{}, error) {
	var err error

	op.O.Raw("SET FOREIGN_KEY_CHECKS=0").Exec()
	for _, table := range dbTables {
		if _, err = op.O.Raw("TRUNCATE TABLE `" + table + "`").Exec(); err != nil {
			return nil, err
		}
	}
	op.O.Raw("SET FOREIGN_KEY_CHECKS=1").Exec()

	// init admin
	op.O.Insert(&User{Name: "system"})

	// init root tree tag
	op.O.Insert(&Tag{Name: ""})
	op.O.Insert(&Tag_rel{TagId: 1, SupTagId: 1, Offset: 0})

	// reset cache
	// ugly hack
	initCache(ctrl.Configure)

	return "reset db done", nil
}
