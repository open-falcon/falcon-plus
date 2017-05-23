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

package falcon_portal

// +-----------+------------------+------+-----+---------+-------+
// | Field     | Type             | Null | Key | Default | Extra |
// +-----------+------------------+------+-----+---------+-------+
// | grp_id    | int(10) unsigned | NO   | MUL | NULL    |       |
// | tpl_id    | int(10) unsigned | NO   | MUL | NULL    |       |
// | bind_user | varchar(64)      | NO   |     |         |       |
// +-----------+------------------+------+-----+---------+-------+

type GrpTpl struct {
	GrpID    int64  `json:"grp_id" gorm:"column:grp_id"`
	TplID    int64  `json:"tpl_id" gorm:"column:tpl_id"`
	BindUser string `json:"bind_user" gorm:"column:bind_user"`
}

func (this GrpTpl) TableName() string {
	return "grp_tpl"
}
