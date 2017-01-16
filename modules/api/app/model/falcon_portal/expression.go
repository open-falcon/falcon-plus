package falcon_portal

// +-------------+------------------+------+-----+---------+----------------+
// | Field       | Type             | Null | Key | Default | Extra          |
// +-------------+------------------+------+-----+---------+----------------+
// | id          | int(10) unsigned | NO   | PRI | NULL    | auto_increment |
// | expression  | varchar(1024)    | NO   |     | NULL    |                |
// | func        | varchar(16)      | NO   |     | all(#1) |                |
// | op          | varchar(8)       | NO   |     |         |                |
// | right_value | varchar(16)      | NO   |     |         |                |
// | max_step    | int(11)          | NO   |     | 1       |                |
// | priority    | tinyint(4)       | NO   |     | 0       |                |
// | note        | varchar(1024)    | NO   |     |         |                |
// | action_id   | int(10) unsigned | NO   |     | 0       |                |
// | create_user | varchar(64)      | NO   |     |         |                |
// | pause       | tinyint(1)       | NO   |     | 0       |                |
// +-------------+------------------+------+-----+---------+----------------+

type Expression struct {
	ID         int64  `json:"id" gorm:"column:id"`
	Expression string `json:"expression" gorm:"column:expression"`
	Func       string `json:"func" gorm:"column:func"`
	Op         string `json:"op" gorm:"column:op"`
	RightValue string `json:"right_value" gorm:"column:right_value"`
	MaxStep    int    `json:"max_step" gorm:"column:max_step"`
	Priority   int    `json:"priority" gorm:"column:priority"`
	Note       string `json:"note" gorm:"column:note"`
	ActionId   int64  `json:"action_id" gorm:"column:action_id"`
	CreateUser string `json:"create_user" gorm:"column:create_user"`
	Pause      int    `json:"pause" gorm:"column:pause"`
}
