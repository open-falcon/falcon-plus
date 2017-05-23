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

package alarm

import "time"

// +--------------+------------------+------+-----+-------------------+-----------------------------+
// | Field        | Type             | Null | Key | Default           | Extra                       |
// +--------------+------------------+------+-----+-------------------+-----------------------------+
// | id           | mediumint(9)     | NO   | PRI | NULL              | auto_increment              |
// | event_caseId | varchar(50)      | YES  | MUL | NULL              |                             |
// | step         | int(10) unsigned | YES  |     | NULL              |                             |
// | cond         | varchar(200)     | NO   |     | NULL              |                             |
// | status       | int(3) unsigned  | YES  |     | 0                 |                             |
// | timestamp    | timestamp        | NO   |     | CURRENT_TIMESTAMP | on update CURRENT_TIMESTAMP |
// +--------------+------------------+------+-----+-------------------+-----------------------------+

type Events struct {
	ID          int64      `json:"id" gorm:"column:id"`
	EventCaseId string     `json:"event_caseId" gorm:"column:event_caseId"`
	Step        int        `json:"step" grom:"step"`
	Cond        string     `json:"cond" grom:"cond"`
	Status      int        `json:"status" grom:"status"`
	Timestamp   *time.Time `json:"timestamp" grom:"timestamp"`
}

func (this Events) TableName() string {
	return "events"
}
