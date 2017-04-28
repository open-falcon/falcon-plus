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
	"testing"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func testTokenInitDb(t *testing.T, o orm.Ormer) (err error) {
	tables := []string{
		"tag_host",
		"tag_rel",
		"tpl_rel",
		"host",
		"token",
		"tag",
		"role",
		"user",
	}
	t.Log("enter testTokenInitDb")
	o.Raw("SET FOREIGN_KEY_CHECKS=0").Exec()
	for _, table := range tables {
		if _, err = o.Raw("TRUNCATE TABLE `" + table + "`").Exec(); err != nil {
			return
		}
	}
	o.Raw("SET FOREIGN_KEY_CHECKS=1").Exec()

	// init admin
	o.Insert(&User{Name: "admin"})

	// init root tree tag
	o.Insert(&Tag{})

	return nil
}

func TestToken(t *testing.T) {
	var items []string
	tag_idx := make(map[string]int64)
	user_idx := make(map[string]int64)
	role_idx := make(map[string]int64)
	token_idx := make(map[string]int64)

	if !test_db_init {
		t.Log("test db not inited, skip test token")
		return
	}
	o := orm.NewOrm()
	err := testTokenInitDb(t, o)
	if err != nil {
		t.Error("init db failed", err)
	}

	schema, _ := NewTagSchema("a,b,c,d,")
	sys, _ := GetUser(1, o)
	op := &Operator{
		User:  sys,
		O:     o,
		Token: SYS_F_A_TOKEN,
	}

	// tag
	items = []string{
		"a=1",
		"a=1,b=1",
		"a=1,b=2",
		"a=1,b=2,c=1",
		"a=1,b=2,c=2",
	}
	for _, item := range items {
		if tag_idx[item], err = op.addTag(&Tag{Name: item}, schema); err != nil {
			t.Error(err)
		}
	}

	// user
	items = []string{
		"u1",
	}
	for _, item := range items {
		if u, err := op.AddUser(&User{Name: item}); err != nil {
			t.Error(err)
		} else {
			user_idx[item] = u.Id
		}
	}

	// role
	items = []string{
		"r1",
		"r2",
		"r3",
		"r4",
	}
	for _, item := range items {
		if role_idx[item], err = op.AddRole(&Role{Name: item}); err != nil {
			t.Error(err)
		}
	}

	// token
	items = []string{
		"token1",
		"token2",
		"token3",
		"token41",
		"token42",
	}
	for _, item := range items {
		if token_idx[item], err = op.AddToken(&Token{Name: item}); err != nil {
			t.Error(err)
		}
	}

	// bind user
	binds := [][3]int64{
		{tag_idx["a=1,b=2"], role_idx["r1"], user_idx["u1"]},
		{tag_idx["a=1,b=2"], role_idx["r2"], user_idx["u1"]},
		{tag_idx["a=1,b=2"], role_idx["r3"], user_idx["u1"]},
		{tag_idx["a=1,b=2"], role_idx["r4"], user_idx["u1"]},
	}
	for _, n := range binds {
		if _, err := addTplRel(op.O, op.User.Id, n[0], n[1], n[2], TPL_REL_T_ACL_USER); err != nil {
			t.Error(err)
		}
	}

	// bind token
	binds = [][3]int64{
		{tag_idx["a=1,b=2"], role_idx["r1"], token_idx["token1"]},
		{tag_idx["a=1"], role_idx["r2"], token_idx["token2"]},
		{tag_idx["a=1,b=2,c=2"], role_idx["r3"], token_idx["token3"]},
		{tag_idx["a=1"], role_idx["r4"], token_idx["token41"]},
		{tag_idx["a=1,b=2"], role_idx["r4"], token_idx["token42"]},
	}
	for _, n := range binds {
		if _, err := addTplRel(op.O, op.User.Id, n[0], n[1], n[2], TPL_REL_T_ACL_TOKEN); err != nil {
			t.Error(err)
		}
	}

	c0 := []struct {
		name     string
		uid      int64
		token_id int64
		want     []int64
		wante    error
	}{
		{name: "case1", uid: user_idx["u1"], token_id: token_idx["token1"], want: []int64{tag_idx["a=1,b=2"], tag_idx["a=1,b=2,c=1"], tag_idx["a=1,b=2,c=2"]}, wante: nil},
		{name: "case2", uid: user_idx["u1"], token_id: token_idx["token2"], want: []int64{tag_idx["a=1,b=2"], tag_idx["a=1,b=2,c=1"], tag_idx["a=1,b=2,c=2"]}, wante: nil},
		{name: "case3", uid: user_idx["u1"], token_id: token_idx["token3"], want: []int64{}, wante: EACCES},
	}
	for _, c := range c0 {
		if tag_id, gote := userHasToken(op.O, c.uid, c.token_id); gote != c.wante {
			t.Errorf("%s userHasToken(%d,%d) = (%d, %v); want (%v)",
				c.name, c.uid, c.token_id,
				tag_id, gote, c.wante)
		}
		if tag_ids, gote := userHasTokenTags(op.O, c.uid, c.token_id); intscmp64(tag_ids, c.want) != 0 || gote != c.wante {
			t.Errorf("%s userHasTokenTags(%d,%d) = (%v, %v); want (%v, %v)",
				c.name, c.uid, c.token_id,
				tag_ids, gote, c.want, c.wante)
		}

	}

	// case1~4
	cases := []struct {
		name     string
		uid      int64
		token_id int64
		tid      int64
		want     int64
		wante    error
	}{
		//case1
		{name: "case1-1", uid: user_idx["u1"], token_id: token_idx["token1"], tid: tag_idx["a=1"], want: 0, wante: EACCES},
		{name: "case1-2", uid: user_idx["u1"], token_id: token_idx["token1"], tid: tag_idx["a=1,b=2"], want: tag_idx["a=1,b=2"], wante: nil},
		{name: "case1-3", uid: user_idx["u1"], token_id: token_idx["token1"], tid: tag_idx["a=1,b=2,c=1"], want: tag_idx["a=1,b=2"], wante: nil},
		{name: "case1-4", uid: user_idx["u1"], token_id: token_idx["token1"], tid: tag_idx["a=1,b=2,c=2"], want: tag_idx["a=1,b=2"], wante: nil},
		//case2
		{name: "case2-1", uid: user_idx["u1"], token_id: token_idx["token2"], tid: tag_idx["a=1"], want: 0, wante: EACCES},
		{name: "case2-2", uid: user_idx["u1"], token_id: token_idx["token2"], tid: tag_idx["a=1,b=2"], want: tag_idx["a=1"], wante: nil},
		{name: "case2-3", uid: user_idx["u1"], token_id: token_idx["token2"], tid: tag_idx["a=1,b=2,c=1"], want: tag_idx["a=1"], wante: nil},
		{name: "case2-4", uid: user_idx["u1"], token_id: token_idx["token2"], tid: tag_idx["a=1,b=2,c=2"], want: tag_idx["a=1"], wante: nil},
		//case3
		{name: "case3-1", uid: user_idx["u1"], token_id: token_idx["token3"], tid: tag_idx["a=1"], want: 0, wante: EACCES},
		{name: "case3-2", uid: user_idx["u1"], token_id: token_idx["token3"], tid: tag_idx["a=1,b=2"], want: 0, wante: EACCES},
		{name: "case3-3", uid: user_idx["u1"], token_id: token_idx["token3"], tid: tag_idx["a=1,b=2,c=1"], want: 0, wante: EACCES},
		{name: "case3-4", uid: user_idx["u1"], token_id: token_idx["token3"], tid: tag_idx["a=1,b=2,c=2"], want: 0, wante: EACCES},
		//case4
		{name: "case4-1", uid: user_idx["u1"], token_id: token_idx["token41"], tid: tag_idx["a=1"], want: 0, wante: EACCES},
		{name: "case4-2", uid: user_idx["u1"], token_id: token_idx["token41"], tid: tag_idx["a=1,b=2"], want: tag_idx["a=1"], wante: nil},
		{name: "case4-3", uid: user_idx["u1"], token_id: token_idx["token41"], tid: tag_idx["a=1,b=2,c=1"], want: tag_idx["a=1"], wante: nil},
		{name: "case4-4", uid: user_idx["u1"], token_id: token_idx["token41"], tid: tag_idx["a=1,b=2,c=2"], want: tag_idx["a=1"], wante: nil},
		{name: "case4-5", uid: user_idx["u1"], token_id: token_idx["token42"], tid: tag_idx["a=1"], want: 0, wante: EACCES},
		{name: "case4-6", uid: user_idx["u1"], token_id: token_idx["token42"], tid: tag_idx["a=1,b=2"], want: tag_idx["a=1,b=2"], wante: nil},
		{name: "case4-7", uid: user_idx["u1"], token_id: token_idx["token42"], tid: tag_idx["a=1,b=2,c=1"], want: tag_idx["a=1,b=2"], wante: nil},
		{name: "case4-8", uid: user_idx["u1"], token_id: token_idx["token42"], tid: tag_idx["a=1,b=2,c=2"], want: tag_idx["a=1,b=2"], wante: nil},
	}
	for _, c := range cases {
		if got, gote := access(op.O, c.uid, c.token_id,
			c.tid); got != c.want || gote != c.wante {
			t.Errorf("%s access(%d,%d,%d) = (%d, %v); want (%d %v)",
				c.name, c.uid, c.token_id, c.tid,
				got, gote, c.want, c.wante)
		}

	}

}
