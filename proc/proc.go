package proc

import (
	"github.com/open-falcon/task/g"
)

func InitProc() {
}

func GetAll() []interface{} {
	ret := make([]interface{}, 0)
	// g.config
	ret = append(ret, g.Config())

	return ret
}
