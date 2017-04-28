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

func testSettingsInitDb(t *testing.T, o orm.Ormer) (err error) {
	t.Log("enter testSettingsInitDb")
	o.Raw("SET FOREIGN_KEY_CHECKS=0").Exec()
	for _, table := range dbTables {
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

func TestPopulate(t *testing.T) {

	if !test_db_init {
		t.Logf("test db not inited, skip test populate\n")
		return
	}

	o := orm.NewOrm()
	sys, _ := GetUser(1, o)
	op := &Operator{
		O:     o,
		User:  sys,
		Token: SYS_F_A_TOKEN | SYS_F_O_TOKEN,
	}

	err := testSettingsInitDb(t, op.O)
	if err != nil {
		t.Error("init db failed", err)
	}

	if _, err := op.Populate(); err != nil {
		t.Error(err)
	}
}
