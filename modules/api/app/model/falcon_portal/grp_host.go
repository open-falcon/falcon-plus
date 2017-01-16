package falcon_portal

// +---------+------------------+------+-----+---------+-------+
// | Field   | Type             | Null | Key | Default | Extra |
// +---------+------------------+------+-----+---------+-------+
// | grp_id  | int(10) unsigned | NO   | PRI | NULL    |       |
// | host_id | int(11)          | NO   | PRI | NULL    |       |
// +---------+------------------+------+-----+---------+-------+

type GrpHost struct {
	GrpID  int64 `json:"grp_id" gorm:"column:grp_id"`
	HostID int64 `json:"host_id" gorm:"column:host_id"`
}

func (this GrpHost) TableName() string {
	return "grp_host"
}
