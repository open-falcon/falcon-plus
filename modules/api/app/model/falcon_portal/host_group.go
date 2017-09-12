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

// +-------------+------------------+------+-----+-------------------+----------------+
// | Field       | Type             | Null | Key | Default           | Extra          |
// +-------------+------------------+------+-----+-------------------+----------------+
// | id          | int(10) unsigned | NO   | PRI | NULL              | auto_increment |
// | grp_name    | varchar(255)     | NO   | UNI |                   |                |
// | create_user | varchar(64)      | NO   |     |                   |                |
// | create_at   | timestamp        | NO   |     | CURRENT_TIMESTAMP |                |
// | come_from   | tinyint(4)       | NO   |     | 0                 |                |
// +-------------+------------------+------+-----+-------------------+----------------+

type HostGroup struct {
	ID         int64  `json:"id" gorm:"column:id"`
	Name       string `json:"grp_name" gorm:"column:grp_name"`
	CreateUser string `json:"create_user" gorm:"column:create_user"`
	ComeFrom   int    `json:"-"  gorm:"column:come_from"`
}

func (this HostGroup) TableName() string {
	return "grp"
}
