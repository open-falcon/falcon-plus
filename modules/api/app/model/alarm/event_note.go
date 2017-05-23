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

import (
	"time"

	"github.com/open-falcon/falcon-plus/modules/api/app/model/uic"
	"github.com/open-falcon/falcon-plus/modules/api/config"
)

// +--------------+------------------+------+-----+-------------------+-----------------------------+
// | Field        | Type             | Null | Key | Default           | Extra                       |
// +--------------+------------------+------+-----+-------------------+-----------------------------+
// | id           | mediumint(9)     | NO   | PRI | NULL              | auto_increment              |
// | event_caseId | varchar(50)      | YES  | MUL | NULL              |                             |
// | note         | varchar(300)     | YES  |     | NULL              |                             |
// | case_id      | varchar(20)      | YES  |     | NULL              |                             |
// | status       | varchar(15)      | YES  |     | NULL              |                             |
// | timestamp    | timestamp        | NO   |     | CURRENT_TIMESTAMP | on update CURRENT_TIMESTAMP |
// | user_id      | int(10) unsigned | YES  | MUL | NULL              |                             |
// +--------------+------------------+------+-----+-------------------+-----------------------------+

type EventNote struct {
	ID          int64      `json:"id" gorm:"column:id"`
	EventCaseId string     `json:"event_caseId" gorm:"column:event_caseId"`
	Note        string     `json:"note" grom:"note"`
	CaseId      string     `json:"case_id" grom:"case_id"`
	Status      string     `json:"status" grom:"status"`
	Timestamp   *time.Time `json:"timestamp" grom:"timestamp"`
	UserId      int64      `json:"user_id" grom:"user_id"`
}

func (this EventNote) TableName() string {
	return "event_note"
}

func (this EventNote) GetUserName() string {
	db := config.Con()
	user := uic.User{ID: this.UserId}
	db.Uic.Table(user.TableName()).Where(&user).Scan(&user)
	return user.Name
}
