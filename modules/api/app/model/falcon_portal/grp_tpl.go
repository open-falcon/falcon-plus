package falcon_portal

// +-----------+------------------+------+-----+---------+-------+
// | Field     | Type             | Null | Key | Default | Extra |
// +-----------+------------------+------+-----+---------+-------+
// | grp_id    | int(10) unsigned | NO   | MUL | NULL    |       |
// | tpl_id    | int(10) unsigned | NO   | MUL | NULL    |       |
// | bind_user | varchar(64)      | NO   |     |         |       |
// +-----------+------------------+------+-----+---------+-------+

type GrpTpl struct {
	GrpID    int64  `json:"grp_id" gorm:"column:grp_id"`
	TplID    int64  `json:"tpl_id" gorm:"column:tpl_id"`
	BindUser string `json:"bind_user" gorm:"column:bind_user"`
}

func (this GrpTpl) TableName() string {
	return "grp_tpl"
}
