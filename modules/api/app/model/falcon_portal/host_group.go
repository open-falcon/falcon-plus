package falcon_portal

// +-------------+------------------+------+-----+-------------------+----------------+
// | Field       | Type             | Null | Key | Default           | Extra          |
// +-------------+------------------+------+-----+-------------------+----------------+
// | id          | int(10) unsigned | NO   | PRI | NULL              | auto_increment |
// | grp_name    | varchar(255)     | NO   | UNI |                   |                |
// | create_user | varchar(64)      | NO   |     |                   |                |
// | create_at   | timestamp        | NO   |     | CURRENT_TIMESTAMP |                |
// | come_from   | tinyint(4)       | NO   |     | 0                 |                |
// +-------------+------------------+------+-----+-------------------+----------------+

type HostGroup struct {
	ID         int64  `json:"id" gorm:"column:id"`
	Name       string `json:"grp_name" gorm:"column:grp_name"`
	CreateUser string `json:"create_user" gorm:"column:create_user"`
	ComeFrom   int    `json:"-"  gorm:"column:come_from"`
}

func (this HostGroup) TableName() string {
	return "grp"
}
