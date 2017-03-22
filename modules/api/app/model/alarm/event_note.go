package alarm

import "time"

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
	ID          int64     `json:"id" gorm:"column:id"`
	EventCaseId string    `json:"event_caseId" grom:"column:event_caseId"`
	Note        string    `json:"note" grom:"note"`
	CaseId      string    `json:"case_id" grom:"case_id"`
	Status      string    `json:"status" grom:"status"`
	Timestamp   time.Time `json:"timestamp" grom:"timestamp"`
	UserId      int64     `json:"user_id" grom:"user_id"`
}

func (this EventNote) TableName() string {
	return "event_noe"
}
