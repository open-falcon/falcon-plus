package falcon_portal

// +-------------+------------------+------+-----+-------------------+----------------+
// | Field       | Type             | Null | Key | Default           | Extra          |
// +-------------+------------------+------+-----+-------------------+----------------+
// | id          | int(10) unsigned | NO   | PRI | NULL              | auto_increment |
// | grp_id      | int(10) unsigned | NO   | MUL | NULL              |                |
// | dir         | varchar(255)     | NO   |     | NULL              |                |
// | create_user | varchar(64)      | NO   |     |                   |                |
// | create_at   | timestamp        | NO   |     | CURRENT_TIMESTAMP |                |
// +-------------+------------------+------+-----+-------------------+----------------+

type Plugin struct {
	ID         int64  `json:"id" gorm:"column:id"`
	GrpId      int64  `json:"grp_id" gorm:"column:grp_id"`
	Dir        string `json:"dir" gorm:"column:dir"`
	CreateUser string `json:"create_user" gorm:"column:create_user"`
}

func (this Plugin) TableName() string {
	return "plugin_dir"
}
