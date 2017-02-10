package dashboard

// +-----------+------------------+------+-----+-------------------+----------------+
// | Field     | Type             | Null | Key | Default           | Extra          |
// +-----------+------------------+------+-----+-------------------+----------------+
// | id        | int(11) unsigned | NO   | PRI | NULL              | auto_increment |
// | endpoints | varchar(10240)   | NO   |     |                   |                |
// | counters  | varchar(10240)   | NO   |     |                   |                |
// | ck        | varchar(32)      | NO   | UNI | NULL              |                |
// | time_     | timestamp        | NO   |     | CURRENT_TIMESTAMP |                |
// +-----------+------------------+------+-----+-------------------+----------------+

type DashboardTmpGraph struct {
	ID        int64  `json:"id" gorm:"column:id"`
	Endpoints string `json:"endpoints" gorm:"column:endpoints"`
	Counters  string `json:"counters" gorm:"column:counters"`
	CK        string `json:"ck" gorm:"column:ck"`
}

func (this DashboardTmpGraph) TableName() string {
	return "tmp_graph"
}
