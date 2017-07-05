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
	Note        string     `json:"note" gorm:"note"`
	CaseId      string     `json:"case_id" gorm:"case_id"`
	Status      string     `json:"status" gorm:"status"`
	Timestamp   *time.Time `json:"timestamp" gorm:"timestamp"`
	UserId      int64      `json:"user_id" gorm:"user_id"`
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
