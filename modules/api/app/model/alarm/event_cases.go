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
	"fmt"
	"time"

	"github.com/open-falcon/falcon-plus/modules/api/config"
)

// +----------------+------------------+------+-----+-------------------+-----------------------------+
// | Field          | Type             | Null | Key | Default           | Extra                       |
// +----------------+------------------+------+-----+-------------------+-----------------------------+
// | id             | varchar(50)      | NO   | PRI | NULL              |                             |
// | endpoint       | varchar(100)     | NO   | MUL | NULL              |                             |
// | metric         | varchar(200)     | NO   |     | NULL              |                             |
// | func           | varchar(50)      | YES  |     | NULL              |                             |
// | cond           | varchar(200)     | NO   |     | NULL              |                             |
// | note           | varchar(500)     | YES  |     | NULL              |                             |
// | max_step       | int(10) unsigned | YES  |     | NULL              |                             |
// | current_step   | int(10) unsigned | YES  |     | NULL              |                             |
// | priority       | int(6)           | NO   |     | NULL              |                             |
// | status         | varchar(20)      | NO   |     | NULL              |                             |
// | timestamp      | timestamp        | NO   |     | CURRENT_TIMESTAMP | on update CURRENT_TIMESTAMP |
// | update_at      | timestamp        | YES  |     | NULL              |                             |
// | closed_at      | timestamp        | YES  |     | NULL              |                             |
// | closed_note    | varchar(250)     | YES  |     | NULL              |                             |
// | user_modified  | int(10) unsigned | YES  |     | NULL              |                             |
// | tpl_creator    | varchar(64)      | YES  |     | NULL              |                             |
// | expression_id  | int(10) unsigned | YES  |     | NULL              |                             |
// | strategy_id    | int(10) unsigned | YES  |     | NULL              |                             |
// | template_id    | int(10) unsigned | YES  |     | NULL              |                             |
// | process_note   | mediumint(9)     | YES  |     | NULL              |                             |
// | process_status | varchar(20)      | YES  |     | unresolved        |                             |
// +----------------+------------------+------+-----+-------------------+-----------------------------+

type EventCases struct {
	ID            string     `json:"id" gorm:"column:id"`
	Endpoint      string     `json:"endpoint" gorm:"column:endpoint"`
	Metric        string     `json:"metric" gorm:"metric"`
	Func          string     `json:"func" gorm:"func"`
	Cond          string     `json:"cond" gorm:"cond"`
	Note          string     `json:"note" gorm:"note"`
	MaxStep       int        `json:"step" gorm:"step"`
	CurrentStep   int        `json:"current_step" gorm:"current_step"`
	Priority      int        `json:"priority" gorm:"priority"`
	Status        string     `json:"status" gorm:"status"`
	Timestamp     *time.Time `json:"timestamp" gorm:"timestamp"`
	UpdateAt      *time.Time `json:"update_at" gorm:"update_at"`
	ClosedAt      *time.Time `json:"closed_at" gorm:"closed_at"`
	ClosedNote    string     `json:"closed_note" gorm:"closed_note"`
	UserModified  int64      `json:"user_modified" gorm:"user_modified"`
	TplCreator    string     `json:"tpl_creator" gorm:"tpl_creator"`
	ExpressionId  int64      `json:"expression_id" gorm:"expression_id"`
	StrategyId    int64      `json:"strategy_id" gorm:"strategy_id"`
	TemplateId    int64      `json:"template_id" gorm:"template_id"`
	ProcessNote   int64      `json:"process_note" gorm:"process_note"`
	ProcessStatus string     `json:"process_status" gorm:"process_status"`
}

func (this EventCases) TableName() string {
	return "event_cases"
}

func (this EventCases) GetEvents() []Events {
	db := config.Con()
	t := Events{
		EventCaseId: this.ID,
	}
	e := []Events{}
	db.Alarm.Table(t.TableName()).Where(&t).Scan(&e)
	return e
}

func (this EventCases) GetNotes() []EventNote {
	db := config.Con()
	perpareSql := fmt.Sprintf("event_caseId = '%s' AND timestamp >= FROM_UNIXTIME(%d)", this.ID, this.Timestamp.Unix())
	t := EventCases{}
	notes := []EventNote{}
	db.Alarm.Table(t.TableName()).Where(perpareSql).Scan(&notes)
	return notes
}

func (this EventCases) NotesCount() int {
	notes := this.GetNotes()
	return len(notes)
}
