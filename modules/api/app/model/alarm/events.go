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
	Step        int        `json:"step" gorm:"step"`
	Cond        string     `json:"cond" gorm:"cond"`
	Status      int        `json:"status" gorm:"status"`
	Timestamp   *time.Time `json:"timestamp" gorm:"timestamp"`
}

func (this Events) TableName() string {
	return "events"
}
