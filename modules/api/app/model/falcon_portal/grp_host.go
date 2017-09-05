package falcon_portal

import (
	con "github.com/open-falcon/falcon-plus/modules/api/config"
)

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

func (this GrpHost) Existing() bool {
	var tGrpHost GrpHost
	db := con.Con()
	db.Falcon.Table(this.TableName()).Where("grp_id = ? AND host_id = ?", this.GrpID, this.HostID).Scan(&tGrpHost)
	if tGrpHost.GrpID != 0 {
		return true
	} else {
		return false
	}
}
