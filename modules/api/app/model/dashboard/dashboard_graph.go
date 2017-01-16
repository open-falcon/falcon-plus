package dashboard

// +-------------+------------------+------+-----+---------+----------------+
// | Field       | Type             | Null | Key | Default | Extra          |
// +-------------+------------------+------+-----+---------+----------------+
// | id          | int(11) unsigned | NO   | PRI | NULL    | auto_increment |
// | title       | char(128)        | NO   |     | NULL    |                |
// | hosts       | varchar(10240)   | NO   |     |         |                |
// | counters    | varchar(1024)    | NO   |     |         |                |
// | screen_id   | int(11) unsigned | NO   | MUL | NULL    |                |
// | timespan    | int(11) unsigned | NO   |     | 3600    |                |
// | graph_type  | char(2)          | NO   |     | h       |                |
// | method      | char(8)          | YES  |     |         |                |
// | position    | int(11) unsigned | NO   |     | 0       |                |
// | falcon_tags | varchar(512)     | NO   |     |         |                |
// +-------------+------------------+------+-----+---------+----------------+

type DashboardGraph struct {
	ID         int64  `json:"id" gorm:"column:id"`
	Title      string `json:"title" gorm:"column:title"`
	Hosts      string `json:"hosts" gorm:"column:hosts"`
	Counters   string `json:"counters" gorm:"column:counters"`
	ScreenId   int64  `json:"screen_id" gorm:"column:screen_id"`
	TimeSpan   int    `json:"timespan" gorm:"column:timespan"`
	GraphType  string `json:"graph_type" gorm:"column:graph_type"`
	Method     string `json:"method" gorm:"column:method"`
	Position   int    `json:"position" gorm:"column:position"`
	FalconTags string `json:"falcon_tags" gorm:"column:falcon_tags"`
}

func (this DashboardGraph) TableName() string {
	return "dashboard_graph"
}
