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

// +-------------+------------------+------+-----+---------+----------------+
// | Field       | Type             | Null | Key | Default | Extra          |
// +-------------+------------------+------+-----+---------+----------------+
// | id          | int(10) unsigned | NO   | PRI | NULL    | auto_increment |
// | expression  | varchar(1024)    | NO   |     | NULL    |                |
// | func        | varchar(16)      | NO   |     | all(#1) |                |
// | op          | varchar(8)       | NO   |     |         |                |
// | right_value | varchar(16)      | NO   |     |         |                |
// | max_step    | int(11)          | NO   |     | 1       |                |
// | priority    | tinyint(4)       | NO   |     | 0       |                |
// | note        | varchar(1024)    | NO   |     |         |                |
// | action_id   | int(10) unsigned | NO   |     | 0       |                |
// | create_user | varchar(64)      | NO   |     |         |                |
// | pause       | tinyint(1)       | NO   |     | 0       |                |
// +-------------+------------------+------+-----+---------+----------------+

type Expression struct {
	ID         int64  `json:"id" gorm:"column:id"`
	Expression string `json:"expression" gorm:"column:expression"`
	Func       string `json:"func" gorm:"column:func"`
	Op         string `json:"op" gorm:"column:op"`
	RightValue string `json:"right_value" gorm:"column:right_value"`
	MaxStep    int    `json:"max_step" gorm:"column:max_step"`
	Priority   int    `json:"priority" gorm:"column:priority"`
	Note       string `json:"note" gorm:"column:note"`
	ActionId   int64  `json:"action_id" gorm:"column:action_id"`
	CreateUser string `json:"create_user" gorm:"column:create_user"`
	Pause      int    `json:"pause" gorm:"column:pause"`
}
