package falcon_portal

import (
	con "github.com/open-falcon/falcon-plus/modules/api/config"
)

// +-------------+------------------+------+-----+-------------------+-----------------------------+
// | Field       | Type             | Null | Key | Default           | Extra                       |
// +-------------+------------------+------+-----+-------------------+-----------------------------+
// | id          | int(10) unsigned | NO   | PRI | NULL              | auto_increment              |
// | grp_id      | int(11)          | NO   |     | NULL              |                             |
// | numerator   | varchar(10240)   | NO   |     | NULL              |                             |
// | denominator | varchar(10240)   | NO   |     | NULL              |                             |
// | endpoint    | varchar(255)     | NO   |     | NULL              |                             |
// | metric      | varchar(255)     | NO   |     | NULL              |                             |
// | tags        | varchar(255)     | NO   |     | NULL              |                             |
// | ds_type     | varchar(255)     | NO   |     | NULL              |                             |
// | step        | int(11)          | NO   |     | NULL              |                             |
// | last_update | timestamp        | NO   |     | CURRENT_TIMESTAMP | on update CURRENT_TIMESTAMP |
// | creator     | varchar(255)     | NO   |     | NULL              |                             |
// +-------------+------------------+------+-----+-------------------+-----------------------------+

type Cluster struct {
	ID          int64  `json:"id" gorm:"column:id"`
	GrpId       int64  `json:"grp_id" gorm:"column:grp_id"`
	Numerator   string `json:"numerator" gorm:"column:numerator"`
	Denominator string `json:"denominator" gorm:"denominator"`
	Endpoint    string `json:"endpoint" gorm:"endpoint"`
	Metric      string `json:"metric" gorm:"metric"`
	Tags        string `json:"tags" gorm:"tags"`
	DsType      string `json:"ds_type" grom:"ds_type"`
	Step        int    `json:"step" gorm:"step"`
	Creator     string `json:"creator" gorm:"creator"`
}

func (this Cluster) TableName() string {
	return "cluster"
}

func (this Cluster) HostGroupName() (name string, err error) {
	if this.GrpId == 0 {
		return
	}
	db := con.Con()
	var hg HostGroup
	hg.ID = this.GrpId
	if dt := db.Falcon.Find(&hg); dt.Error != nil {
		return name, dt.Error
	}
	name = hg.Name
	return
}
