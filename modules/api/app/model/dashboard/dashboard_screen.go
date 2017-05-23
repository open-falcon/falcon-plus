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

package dashboard

// +-------+------------------+------+-----+-------------------+-----------------------------+
// | Field | Type             | Null | Key | Default           | Extra                       |
// +-------+------------------+------+-----+-------------------+-----------------------------+
// | id    | int(11) unsigned | NO   | PRI | NULL              | auto_increment              |
// | pid   | int(11) unsigned | NO   | MUL | 0                 |                             |
// | name  | char(128)        | NO   |     | NULL              |                             |
// | time  | timestamp        | NO   |     | CURRENT_TIMESTAMP | on update CURRENT_TIMESTAMP |
// +-------+------------------+------+-----+-------------------+-----------------------------+

type DashboardScreen struct {
	ID   int64  `json:"id" gorm:"column:id"`
	PID  int64  `json:"pid" gorm:"column:pid"`
	Name string `json:"name" gorm:"column:name"`
}

func (this DashboardScreen) TableName() string {
	return "dashboard_screen"
}
