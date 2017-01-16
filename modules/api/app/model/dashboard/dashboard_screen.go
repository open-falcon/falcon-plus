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
