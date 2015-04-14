package proc

import (
	P "github.com/open-falcon/model/proc"
	"github.com/open-falcon/task/g"
	"time"
)

// 索引垃圾清理
var (
	IndexDeleteCnt = P.NewSCounterBase("IndexDeleteCnt")
)

func InitProc() {
}

func GetAll() []interface{} {
	ret := make([]interface{}, 0)

	// g.config
	ret = append(ret, g.Config())

	// index
	ret = append(ret, IndexDeleteCnt)

	return ret
}

// TODO 临时放在这里了, 考虑放到合适的模块
func FmtUnixTs(ts int64) string {
	return time.Unix(ts, 0).Format("2006-01-02 15:04:05")
}
