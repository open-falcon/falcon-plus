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
	Endpoint      string     `json:"endpoint" grom:"column:endpoint"`
	Metric        string     `json:"metric" grom:"metric"`
	Func          string     `json:"func" grom:"func"`
	Cond          string     `json:"cond" grom:"cond"`
	Note          string     `json:"note" grom:"note"`
	MaxStep       int        `json:"step" grom:"step"`
	CurrentStep   int        `json:"current_step" grom:"current_step"`
	Priority      int        `json:"priority" grom:"priority"`
	Status        string     `json:"status" grom:"status"`
	Timestamp     *time.Time `json:"timestamp" grom:"timestamp"`
	UpdateAt      *time.Time `json:"update_at" grom:"update_at"`
	ClosedAt      *time.Time `json:"closed_at" grom:"closed_at"`
	ClosedNote    string     `json:"closed_note" grom:"closed_note"`
	UserModified  int64      `json:"user_modified" grom:"user_modified"`
	TplCreator    string     `json:"tpl_creator" grom:"tpl_creator"`
	ExpressionId  int64      `json:"expression_id" grom:"expression_id"`
	StrategyId    int64      `json:"strategy_id" grom:"strategy_id"`
	TemplateId    int64      `json:"template_id" grom:"template_id"`
	ProcessNote   int64      `json:"process_note" grom:"process_note"`
	ProcessStatus string     `json:"process_status" grom:"process_status"`
}

func (this EventCases) TableName() string {
	return "event_cases"
}

var db = config.Con()

func (this EventCases) GetEvents() []Events {
	t := Events{
		EventCaseId: this.ID,
	}
	e := []Events{}
	db.Alarm.Table(t.TableName()).Where(&t).Scan(&e)
	return e
}

func (this EventCases) GetNotes() []EventNote {
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
