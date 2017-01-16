package falcon_portal

// +----------+---------------------+------+-----+-------------------+-----------------------------+
// | Field    | Type                | Null | Key | Default           | Extra                       |
// +----------+---------------------+------+-----+-------------------+-----------------------------+
// | id       | bigint(20) unsigned | NO   | PRI | NULL              | auto_increment              |
// | name     | varchar(255)        | NO   | UNI |                   |                             |
// | obj      | varchar(10240)      | NO   |     |                   |                             |
// | obj_type | varchar(255)        | NO   |     |                   |                             |
// | metric   | varchar(128)        | NO   |     |                   |                             |
// | tags     | varchar(1024)       | NO   |     |                   |                             |
// | dstype   | varchar(32)         | NO   |     | GAUGE             |                             |
// | step     | int(11) unsigned    | NO   |     | 60                |                             |
// | mock     | double              | NO   |     | 0                 |                             |
// | creator  | varchar(64)         | NO   |     |                   |                             |
// | t_create | datetime            | NO   |     | NULL              |                             |
// | t_modify | timestamp           | NO   |     | CURRENT_TIMESTAMP | on update CURRENT_TIMESTAMP |
// +----------+---------------------+------+-----+-------------------+-----------------------------+

//no_data
type Mockcfg struct {
	ID   int64  `json:"id" gorm:"column:id"`
	Name string `json:"name" gorm:"column:name"`
	Obj  string `json:"obj" gorm:"column:obj"`
	//group, host, other
	ObjType string  `json:"obj_type" gorm:"column:obj_type"`
	Metric  string  `json:"metric" gorm:"column:metric"`
	Tags    string  `json:"tags" gorm:"column:tags"`
	DsType  string  `json:"dstype" gorm:"column:dstype"`
	Step    int     `json:"step" gorm:"column:step"`
	Mock    float64 `json:"mock" gorm:"column:mock"`
	Creator string  `json:"creator" gorm:"column:creator"`
}

func (this Mockcfg) TableName() string {
	return "mockcfg"
}
